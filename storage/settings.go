package storage

import "vpub/model"

func (s *Storage) Settings() (model.Settings, error) {
	var settings model.Settings

	err := s.db.QueryRow(`
        SELECT
            name, css, footer, per_page, url, lang
        FROM
            settings;
    `).Scan(
		&settings.Name,
		&settings.Css,
		&settings.Footer,
		&settings.PerPage,
		&settings.URL,
		&settings.Lang,
	)

	return settings, err
}

func (s *Storage) UpdateSettings(settings model.Settings) error {
	query := `
        UPDATE settings 
        SET name=$1, 
            css=$2,
            footer=$3,
            per_page=$4,
            url=$5,
            lang=$6;
    `

	_, err := s.db.Exec(
		query,
		settings.Name,
		settings.Css,
		settings.Footer,
		&settings.PerPage,
		&settings.URL,
		&settings.Lang,
	)

	return err
}
