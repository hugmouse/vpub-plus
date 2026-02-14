-- Add settings cache TTL
ALTER TABLE settings
ADD COLUMN settings_cache_ttl integer NOT NULL DEFAULT 30;
