-- Inisialisasi database fleetdb
CREATE TABLE IF NOT EXISTS vehicle_locations (
    id SERIAL PRIMARY KEY,
    vehicle_id VARCHAR(20) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp BIGINT NOT NULL
);
