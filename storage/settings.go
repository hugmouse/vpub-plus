package storage

import (
	"time"
	"vpub/model"
)

func (s *Storage) Settings() (model.Settings, error) {
	s.settingsMu.RLock()
	if s.settingsCache != nil && time.Now().Before(s.settingsCacheTTL) {
		settings := *s.settingsCache
		s.settingsMu.RUnlock()
		return settings, nil
	}
	s.settingsMu.RUnlock()

	var settings model.Settings

	err := s.db.QueryRow(`
        SELECT
            name, css, footer, per_page, url, lang, image_proxy_cache_time, image_proxy_size_limit, settings_cache_ttl
        FROM
            settings
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
	)

	if err != nil {
		return settings, err
	}

	ttl := time.Duration(settings.SettingsCacheTTL) * time.Second
	if ttl <= 0 {
		ttl = 30 * time.Second
	}

	s.settingsMu.Lock()
	s.settingsCache = &settings
	s.settingsCacheTTL = time.Now().Add(ttl)
	s.settingsMu.Unlock()

	return settings, nil
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
			image_proxy_size_limit=$8,
			settings_cache_ttl=$9;
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
		&settings.SettingsCacheTTL,
	)

	if err != nil {
		return err
	}

	s.settingsMu.Lock()
	s.settingsCache = nil
	s.settingsCacheTTL = time.Time{}
	s.settingsMu.Unlock()

	return nil
}
