package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vpub/model"
)

// SetupRequest bundles the data collected during initial onboarding so that
// CompleteSetup can apply it in a single transaction
type SetupRequest struct {
	AdminName     string
	AdminPassword string
	ForumName     string
	URL           string
	Lang          string
}

// ErrSetupCompleted is returned by CompleteSetup when an admin user already
// exists, for example, onboarding has already finished
type ErrSetupCompleted struct{}

func (ErrSetupCompleted) Error() string {
	return "setup already completed"
}

// welcomeContent is the body of the topic that is seeded on first install
const welcomeContent = `Welcome! You just finished the initial setup.

## What now?

- Head to the [admin area](/admin) to manage forums, boards, users and instance settings.
- Create new boards under **Getting started** or remove this sample content entirely.
- Adjust the forum name, URL and appearance under [instance settings](/admin/settings/edit).

Enjoy your forum!`

// CompleteSetup performs the full initial onboarding — creating the admin user,
// updating the instance settings and seeding sample content — in a single
// transaction. If any step fails the whole install is rolled back, leaving the
// database as if setup never ran, so the operator can retry. It returns the new
// admin user's id and the resulting settings.
//
// The admin-exists check inside the transaction makes CompleteSetup a one-time
// operation: a concurrent caller that also passed the earlier HasAdmin check
// will get ErrSetupCompleted instead of creating a second admin.
func (s *Storage) CompleteSetup(req SetupRequest) (int64, model.Settings, error) {
	var (
		adminID  int64
		settings model.Settings
	)

	hash, err := hashPassword(req.AdminPassword)
	if err != nil {
		return adminID, settings, fmt.Errorf("failed to hash password: %w", err)
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return adminID, settings, err
	}

	// Refuse if an admin already exists.
	var adminExists bool
	if err := tx.QueryRowContext(ctx, `SELECT true FROM users WHERE is_admin=true LIMIT 1`).Scan(&adminExists); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return adminID, settings, rollbackTx(tx, err)
	}
	if adminExists {
		return adminID, settings, rollbackTx(tx, ErrSetupCompleted{})
	}

	// Refuse if the chosen admin name is taken.
	var userExists bool
	if err := tx.QueryRowContext(ctx, `SELECT true FROM users WHERE name=lower($1) LIMIT 1`, req.AdminName).Scan(&userExists); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return adminID, settings, rollbackTx(tx, err)
	}
	if userExists {
		return adminID, settings, rollbackTx(tx, ErrUserExists{})
	}

	if err := tx.QueryRowContext(ctx, `insert into users (name, hash, is_admin) values (lower($1), $2, $3) returning id`, req.AdminName, string(hash), true).Scan(&adminID); err != nil {
		return adminID, settings, rollbackTx(tx, err)
	}

	// Load the existing settings row, apply the onboarding values and persist it.
	if err := tx.QueryRowContext(ctx, `
        SELECT
            name, css, footer, per_page, url, lang, image_proxy_cache_time, image_proxy_size_limit, settings_cache_ttl
        FROM settings
        LIMIT 1;
    `).Scan(
		&settings.Name,
		&settings.CSS,
		&settings.Footer,
		&settings.PerPage,
		&settings.URL,
		&settings.Lang,
		&settings.ImageProxyCacheTime,
		&settings.ImageProxySizeLimit,
		&settings.SettingsCacheTTL,
	); err != nil {
		return adminID, settings, rollbackTx(tx, err)
	}
	settings.Name = req.ForumName
	settings.URL = req.URL
	settings.Lang = req.Lang
	if _, err := tx.ExecContext(ctx, `
        UPDATE settings
        SET name=$1,
            css=$2,
            footer=$3,
            per_page=$4,
            url=$5,
            lang=$6,
            image_proxy_cache_time=$7,
            image_proxy_size_limit=$8,
            settings_cache_ttl=$9;
    `,
		settings.Name,
		settings.CSS,
		settings.Footer,
		settings.PerPage,
		settings.URL,
		settings.Lang,
		settings.ImageProxyCacheTime,
		settings.ImageProxySizeLimit,
		settings.SettingsCacheTTL,
	); err != nil {
		return adminID, settings, rollbackTx(tx, err)
	}

	// Seed sample content: a welcome forum, board and topic.
	var forumID int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO forums (name, position, is_locked) VALUES ($1, $2, $3) RETURNING id`, "Getting started", 0, false).Scan(&forumID); err != nil {
		return adminID, settings, rollbackTx(tx, fmt.Errorf("unable to create forum: %w", err))
	}

	var boardID int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO boards (name, description, position, forum_id, is_locked) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Your first board", "This is a sample board. Feel free to edit or remove it from the admin area.", 0, forumID, false,
	).Scan(&boardID); err != nil {
		return adminID, settings, rollbackTx(tx, fmt.Errorf("unable to create board: %w", err))
	}

	var topicID int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (is_sticky, is_locked, board_id, post_id) VALUES ($1, $2, $3, -1) RETURNING id`, true, true, boardID).Scan(&topicID); err != nil {
		return adminID, settings, rollbackTx(tx, fmt.Errorf("unable to create topic: %w", err))
	}
	var postID int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (subject, content, topic_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		"Welcome to your new forum", welcomeContent, topicID, adminID,
	).Scan(&postID); err != nil {
		return adminID, settings, rollbackTx(tx, fmt.Errorf("unable to create post: %w", err))
	}
	if _, err := tx.ExecContext(ctx, `UPDATE topics set post_id=$1 where id=$2`, postID, topicID); err != nil {
		return adminID, settings, rollbackTx(tx, fmt.Errorf("unable to link topic: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return adminID, settings, err
	}

	s.settingsMu.Lock()
	s.settingsCache = nil
	s.settingsCacheTTL = time.Time{}
	s.settingsMu.Unlock()

	return adminID, settings, nil
}

// rollbackTx rolls tx back and joins any rollback error with the original reason
func rollbackTx(tx *sql.Tx, cause error) error {
	if rbErr := tx.Rollback(); rbErr != nil {
		return errors.Join(cause, fmt.Errorf("rollback in CompleteSetup failed: %w", rbErr))
	}
	return cause
}
