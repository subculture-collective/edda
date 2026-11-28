-- +goose Up
-- Server-owned SRD 5.2.1 route facts (#144): distance in feet, Travel Terrain
-- entry and good road. NULL means undeclared; travel_time prose is not read.
ALTER TABLE location_connections ADD COLUMN srd521_route JSONB;
-- +goose Down
ALTER TABLE location_connections DROP COLUMN IF EXISTS srd521_route;
