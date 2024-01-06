-- name: AddProperty :exec
INSERT INTO properties(
        id,
        humanReadableId,
        created_at,
        updated_at,
        title,
        floor,
        user_id,
        images,
        thumbnail,
        is_featured,
        energy_class,
        energy_consumption_primary,
        energy_emissions_index,
        energy_consumption_green,
        destination_residential,
        destination_commercial,
        destination_office,
        destination_holiday,
        other_utilities_terrance,
        other_utilities_service_toilet,
        other_utilities_underground_storage,
        other_utilities_storage,
        property_transaction,
        furnished_not,
        furnished_partially,
        furnished_complete,
        furnished_luxury,
        interior_needs_renovation,
        interior_has_renovation,
        interior_good_state,
        heating_termoficare,
        heating_central_heating,
        heating_building,
        heating_stove,
        heating_radiator,
        heating_other_electrical,
        heating_gas_convector,
        heating_infrared_panels,
        heating_floor_heating
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14,
        $15,
        $16,
        $17,
        $18,
        $19,
        $20,
        $21,
        $22,
        $23,
        $24,
        $25,
        $26,
        $27,
        $28,
        $29,
        $30,
        $31,
        $32,
        $33,
        $34,
        $35,
        $36,
        $37,
        $38,
        $39
    );
-- name: GetProperty :one
SELECT *
FROM properties
WHERE id = $1
LIMIT 1;
-- name: GetByHumanReadableId :one
SELECT *
FROM properties
WHERE humanReadableId = $1
LIMIT 1;
-- name: DeletePropertyByHumanReadableId :exec
DELETE FROM properties
WHERE humanReadableId = $1;
-- name: ListProperties :many
SELECT *
FROM properties
ORDER BY created_at DESC;