package storage

import "vpub/model"

func (s *Storage) Settings() (model.Settings, error) {
	var settings model.Settings
	err := s.db.QueryRow(`
        SELECT
            name, css
        FROM
            settings;
    `).Scan(&settings.Name, &settings.Css)
	return settings, err
}

func (s *Storage) UpdateSettings(settings model.Settings) error {
	stmt, err := s.db.Prepare(`UPDATE settings SET name=$1, css=$2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(settings.Name, settings.Css)
	return err
}
