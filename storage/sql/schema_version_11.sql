-- Add image proxy settings
ALTER TABLE settings
ADD COLUMN image_proxy_cache_time integer NOT NULL DEFAULT 600,
ADD COLUMN image_proxy_size_limit integer NOT NULL DEFAULT 524288;
