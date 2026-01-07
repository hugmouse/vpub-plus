package storage

import "vpub/model"

func (s *Storage) Settings() (model.Settings, error) {
	var settings model.Settings

	err := s.db.QueryRow(`
        SELECT
            name, css, footer, per_page, url, lang, image_proxy_cache_time, image_proxy_size_limit
        FROM
            settings;
    `).Scan(
		&settings.Name,
		&settings.CSS,
		&settings.Footer,
		&settings.PerPage,
		&settings.URL,
		&settings.Lang,
		&settings.ImageProxyCacheTime,
		&settings.ImageProxySizeLimit,
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
            lang=$6,
			image_proxy_cache_time=$7,
			image_proxy_size_limit=$8;
    `

	_, err := s.db.Exec(
		query,
		settings.Name,
		settings.CSS,
		settings.Footer,
		&settings.PerPage,
		&settings.URL,
		&settings.Lang,
		&settings.ImageProxyCacheTime,
		&settings.ImageProxySizeLimit,
	)

	return err
}
