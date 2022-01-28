package storage

import "vpub/model"

func (s *Storage) Settings() (model.Settings, error) {
	var settings model.Settings
	err := s.db.QueryRow(`
        SELECT
            name, css, footer, per_page
        FROM
            settings;
    `).Scan(&settings.Name, &settings.Css, &settings.Footer, &settings.PerPage)
	return settings, err
}

func (s *Storage) UpdateSettings(settings model.Settings) error {
	stmt, err := s.db.Prepare(`UPDATE settings SET name=$1, css=$2, footer=$3, per_page=$4;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(settings.Name, settings.Css, settings.Footer, &settings.PerPage)
	return err
}
