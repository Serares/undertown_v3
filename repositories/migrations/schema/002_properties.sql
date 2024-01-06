-- +goose Up
CREATE TYPE transaction_type AS ENUM ('sell', 'rent');
CREATE TABLE properties (
    id UUID PRIMARY KEY,
    humanReadableId VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    floor INT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    images TEXT [] NOT NULL,
    thumbnail TEXT NOT NULL,
    is_featured BOOLEAN NOT NULL,
    energy_class VARCHAR(255) NOT NULL,
    energy_consumption_primary VARCHAR(255) NOT NULL,
    energy_emissions_index VARCHAR(255) NOT NULL,
    energy_consumption_green VARCHAR(255) NOT NULL,
    -- regenerable sources
    destination_residential BOOLEAN NOT NULL,
    destination_commercial BOOLEAN NOT NULL,
    destination_office BOOLEAN NOT NULL,
    destination_holiday BOOLEAN NOT NULL,
    other_utilities_terrance BOOLEAN NOT NULL,
    -- terasa
    other_utilities_service_toilet BOOLEAN NOT NULL,
    other_utilities_underground_storage BOOLEAN NOT NULL,
    -- boxa subsol
    other_utilities_storage BOOLEAN NOT NULL,
    property_transaction transaction_type NOT NULL,
    -- debara
    furnished_not BOOLEAN NOT NULL,
    furnished_partially BOOLEAN NOT NULL,
    furnished_complete BOOLEAN NOT NULL,
    furnished_luxury BOOLEAN NOT NULL,
    interior_needs_renovation BOOLEAN NOT NULL,
    interior_has_renovation BOOLEAN NOT NULL,
    interior_good_state BOOLEAN NOT NULL,
    heating_termoficare BOOLEAN NOT NULL,
    -- termoficare
    heating_central_heating BOOLEAN NOT NULL,
    heating_building BOOLEAN NOT NULL,
    heating_stove BOOLEAN NOT NULL,
    heating_radiator BOOLEAN NOT NULL,
    heating_other_electrical BOOLEAN NOT NULL,
    heating_gas_convector BOOLEAN NOT NULL,
    heating_infrared_panels BOOLEAN NOT NULL,
    heating_floor_heating BOOLEAN NOT NULL
);
-- +goose Down
DROP TABLE properties;
DROP TYPE transaction_type;