-- SQLite schema adapted from the provided PostgreSQL schema
-- SQLite doesn't support ENUM, so use TEXT and enforce data integrity in your application logic
-- SQLite also doesn't support UUID natively, use TEXT for UUIDs
-- Replace TEXT [] with a TEXT field for the images
-- BOOLEAN fields are stored as INTEGER (0 or 1) in SQLite
CREATE TABLE properties (
    id TEXT PRIMARY KEY,
    humanReadableId TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    floor INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    -- Assuming users table exists with id as TEXT
    FOREIGN KEY (user_id) REFERENCES users(id),
    images TEXT NOT NULL,
    -- Store as a comma-separated string or JSON string
    thumbnail TEXT NOT NULL,
    is_featured INTEGER NOT NULL,
    energy_class TEXT NOT NULL,
    energy_consumption_primary TEXT NOT NULL,
    energy_emissions_index TEXT NOT NULL,
    energy_consumption_green TEXT NOT NULL,
    destination_residential INTEGER NOT NULL,
    destination_commercial INTEGER NOT NULL,
    destination_office INTEGER NOT NULL,
    destination_holiday INTEGER NOT NULL,
    other_utilities_terrance INTEGER NOT NULL,
    other_utilities_service_toilet INTEGER NOT NULL,
    other_utilities_underground_storage INTEGER NOT NULL,
    other_utilities_storage INTEGER NOT NULL,
    property_transaction TEXT CHECK(property_transaction IN ('sell', 'rent')) NOT NULL,
    furnished_not INTEGER NOT NULL,
    furnished_partially INTEGER NOT NULL,
    furnished_complete INTEGER NOT NULL,
    furnished_luxury INTEGER NOT NULL,
    interior_needs_renovation INTEGER NOT NULL,
    interior_has_renovation INTEGER NOT NULL,
    interior_good_state INTEGER NOT NULL,
    heating_termoficare INTEGER NOT NULL,
    heating_central_heating INTEGER NOT NULL,
    heating_building INTEGER NOT NULL,
    heating_stove INTEGER NOT NULL,
    heating_radiator INTEGER NOT NULL,
    heating_other_electrical INTEGER NOT NULL,
    heating_gas_convector INTEGER NOT NULL,
    heating_infrared_panels INTEGER NOT NULL,
    heating_floor_heating INTEGER NOT NULL
);
-- SQLite doesn't support the DROP TYPE statement as it does not support custom types
-- The DROP TABLE statement remains the same
-- +goose Down
DROP TABLE properties;