package storage

import (
	"pboard/model"
)

func (s *Storage) NotificationsByUser(user string) ([]model.Notification, error) {
	q := `
		SELECT
    		n.id, r.id, r.author, r.content, r.post_id, r.parent_id, p.title from notifications n
		LEFT JOIN replies r on n.reply_id = r.id
		LEFT JOIN posts p on r.post_id = p.id
		WHERE n.author=$1
		ORDER BY r.id DESC`
	rows, err := s.db.Query(q, user)
	if err != nil {
		return nil, err
	}
	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		err := rows.Scan(&notification.Id, &notification.Reply.Id, &notification.Reply.Author, &notification.Reply.Content, &notification.Reply.PostId, &notification.Reply.ParentId, &notification.Reply.PostTitle)
		if err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (s *Storage) UserHasNotifications(name string) bool {
	var rv bool
	s.db.QueryRow(`SELECT true FROM notifications WHERE author=$1`, name).Scan(&rv)
	return rv
}

func (s *Storage) DeleteNotification(id int64) error {
	stmt, err := s.db.Prepare(`DELETE from notifications WHERE id = $1;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}

func (s *Storage) DeleteNotificationByReplyId(id int64) error {
	stmt, err := s.db.Prepare(`DELETE from notifications WHERE reply_id = $1;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}
