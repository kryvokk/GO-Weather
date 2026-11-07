CREATE TABLE IF NOT EXISTS cities (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    country TEXT NOT NULL DEFAULT '',
    latitude DOUBLE PRECISION NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude DOUBLE PRECISION NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    timezone TEXT NOT NULL DEFAULT '',
    UNIQUE (latitude, longitude)
);
CREATE TABLE IF NOT EXISTS favorite_cities (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (city_id)
);
CREATE TABLE IF NOT EXISTS user_settings (
    id BIGINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    temperature_unit TEXT NOT NULL DEFAULT 'C' CHECK (temperature_unit IN ('C','F')),
    wind_speed_unit TEXT NOT NULL DEFAULT 'kmh',
    pressure_unit TEXT NOT NULL DEFAULT 'hPa',
    theme TEXT NOT NULL DEFAULT 'dark',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO user_settings (id) VALUES (1) ON CONFLICT (id) DO NOTHING;
