// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: properties.sql

package psql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addProperty = `-- name: AddProperty :exec
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
    )
`

type AddPropertyParams struct {
	ID                               uuid.UUID
	Humanreadableid                  string
	CreatedAt                        time.Time
	UpdatedAt                        time.Time
	Title                            string
	Floor                            int32
	UserID                           uuid.UUID
	Images                           []string
	Thumbnail                        string
	IsFeatured                       bool
	EnergyClass                      string
	EnergyConsumptionPrimary         string
	EnergyEmissionsIndex             string
	EnergyConsumptionGreen           string
	DestinationResidential           bool
	DestinationCommercial            bool
	DestinationOffice                bool
	DestinationHoliday               bool
	OtherUtilitiesTerrance           bool
	OtherUtilitiesServiceToilet      bool
	OtherUtilitiesUndergroundStorage bool
	OtherUtilitiesStorage            bool
	PropertyTransaction              TransactionType
	FurnishedNot                     bool
	FurnishedPartially               bool
	FurnishedComplete                bool
	FurnishedLuxury                  bool
	InteriorNeedsRenovation          bool
	InteriorHasRenovation            bool
	InteriorGoodState                bool
	HeatingTermoficare               bool
	HeatingCentralHeating            bool
	HeatingBuilding                  bool
	HeatingStove                     bool
	HeatingRadiator                  bool
	HeatingOtherElectrical           bool
	HeatingGasConvector              bool
	HeatingInfraredPanels            bool
	HeatingFloorHeating              bool
}

func (q *Queries) AddProperty(ctx context.Context, arg AddPropertyParams) error {
	_, err := q.db.ExecContext(ctx, addProperty,
		arg.ID,
		arg.Humanreadableid,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Floor,
		arg.UserID,
		pq.Array(arg.Images),
		arg.Thumbnail,
		arg.IsFeatured,
		arg.EnergyClass,
		arg.EnergyConsumptionPrimary,
		arg.EnergyEmissionsIndex,
		arg.EnergyConsumptionGreen,
		arg.DestinationResidential,
		arg.DestinationCommercial,
		arg.DestinationOffice,
		arg.DestinationHoliday,
		arg.OtherUtilitiesTerrance,
		arg.OtherUtilitiesServiceToilet,
		arg.OtherUtilitiesUndergroundStorage,
		arg.OtherUtilitiesStorage,
		arg.PropertyTransaction,
		arg.FurnishedNot,
		arg.FurnishedPartially,
		arg.FurnishedComplete,
		arg.FurnishedLuxury,
		arg.InteriorNeedsRenovation,
		arg.InteriorHasRenovation,
		arg.InteriorGoodState,
		arg.HeatingTermoficare,
		arg.HeatingCentralHeating,
		arg.HeatingBuilding,
		arg.HeatingStove,
		arg.HeatingRadiator,
		arg.HeatingOtherElectrical,
		arg.HeatingGasConvector,
		arg.HeatingInfraredPanels,
		arg.HeatingFloorHeating,
	)
	return err
}

const deletePropertyByHumanReadableId = `-- name: DeletePropertyByHumanReadableId :exec
DELETE FROM properties
WHERE humanReadableId = $1
`

func (q *Queries) DeletePropertyByHumanReadableId(ctx context.Context, humanreadableid string) error {
	_, err := q.db.ExecContext(ctx, deletePropertyByHumanReadableId, humanreadableid)
	return err
}

const getByHumanReadableId = `-- name: GetByHumanReadableId :one
SELECT id, humanreadableid, created_at, updated_at, title, floor, user_id, images, thumbnail, is_featured, energy_class, energy_consumption_primary, energy_emissions_index, energy_consumption_green, destination_residential, destination_commercial, destination_office, destination_holiday, other_utilities_terrance, other_utilities_service_toilet, other_utilities_underground_storage, other_utilities_storage, property_transaction, furnished_not, furnished_partially, furnished_complete, furnished_luxury, interior_needs_renovation, interior_has_renovation, interior_good_state, heating_termoficare, heating_central_heating, heating_building, heating_stove, heating_radiator, heating_other_electrical, heating_gas_convector, heating_infrared_panels, heating_floor_heating
FROM properties
WHERE humanReadableId = $1
LIMIT 1
`

func (q *Queries) GetByHumanReadableId(ctx context.Context, humanreadableid string) (Property, error) {
	row := q.db.QueryRowContext(ctx, getByHumanReadableId, humanreadableid)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Humanreadableid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Floor,
		&i.UserID,
		pq.Array(&i.Images),
		&i.Thumbnail,
		&i.IsFeatured,
		&i.EnergyClass,
		&i.EnergyConsumptionPrimary,
		&i.EnergyEmissionsIndex,
		&i.EnergyConsumptionGreen,
		&i.DestinationResidential,
		&i.DestinationCommercial,
		&i.DestinationOffice,
		&i.DestinationHoliday,
		&i.OtherUtilitiesTerrance,
		&i.OtherUtilitiesServiceToilet,
		&i.OtherUtilitiesUndergroundStorage,
		&i.OtherUtilitiesStorage,
		&i.PropertyTransaction,
		&i.FurnishedNot,
		&i.FurnishedPartially,
		&i.FurnishedComplete,
		&i.FurnishedLuxury,
		&i.InteriorNeedsRenovation,
		&i.InteriorHasRenovation,
		&i.InteriorGoodState,
		&i.HeatingTermoficare,
		&i.HeatingCentralHeating,
		&i.HeatingBuilding,
		&i.HeatingStove,
		&i.HeatingRadiator,
		&i.HeatingOtherElectrical,
		&i.HeatingGasConvector,
		&i.HeatingInfraredPanels,
		&i.HeatingFloorHeating,
	)
	return i, err
}

const getProperty = `-- name: GetProperty :one
SELECT id, humanreadableid, created_at, updated_at, title, floor, user_id, images, thumbnail, is_featured, energy_class, energy_consumption_primary, energy_emissions_index, energy_consumption_green, destination_residential, destination_commercial, destination_office, destination_holiday, other_utilities_terrance, other_utilities_service_toilet, other_utilities_underground_storage, other_utilities_storage, property_transaction, furnished_not, furnished_partially, furnished_complete, furnished_luxury, interior_needs_renovation, interior_has_renovation, interior_good_state, heating_termoficare, heating_central_heating, heating_building, heating_stove, heating_radiator, heating_other_electrical, heating_gas_convector, heating_infrared_panels, heating_floor_heating
FROM properties
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetProperty(ctx context.Context, id uuid.UUID) (Property, error) {
	row := q.db.QueryRowContext(ctx, getProperty, id)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Humanreadableid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Floor,
		&i.UserID,
		pq.Array(&i.Images),
		&i.Thumbnail,
		&i.IsFeatured,
		&i.EnergyClass,
		&i.EnergyConsumptionPrimary,
		&i.EnergyEmissionsIndex,
		&i.EnergyConsumptionGreen,
		&i.DestinationResidential,
		&i.DestinationCommercial,
		&i.DestinationOffice,
		&i.DestinationHoliday,
		&i.OtherUtilitiesTerrance,
		&i.OtherUtilitiesServiceToilet,
		&i.OtherUtilitiesUndergroundStorage,
		&i.OtherUtilitiesStorage,
		&i.PropertyTransaction,
		&i.FurnishedNot,
		&i.FurnishedPartially,
		&i.FurnishedComplete,
		&i.FurnishedLuxury,
		&i.InteriorNeedsRenovation,
		&i.InteriorHasRenovation,
		&i.InteriorGoodState,
		&i.HeatingTermoficare,
		&i.HeatingCentralHeating,
		&i.HeatingBuilding,
		&i.HeatingStove,
		&i.HeatingRadiator,
		&i.HeatingOtherElectrical,
		&i.HeatingGasConvector,
		&i.HeatingInfraredPanels,
		&i.HeatingFloorHeating,
	)
	return i, err
}

const listProperties = `-- name: ListProperties :many
SELECT id, humanreadableid, created_at, updated_at, title, floor, user_id, images, thumbnail, is_featured, energy_class, energy_consumption_primary, energy_emissions_index, energy_consumption_green, destination_residential, destination_commercial, destination_office, destination_holiday, other_utilities_terrance, other_utilities_service_toilet, other_utilities_underground_storage, other_utilities_storage, property_transaction, furnished_not, furnished_partially, furnished_complete, furnished_luxury, interior_needs_renovation, interior_has_renovation, interior_good_state, heating_termoficare, heating_central_heating, heating_building, heating_stove, heating_radiator, heating_other_electrical, heating_gas_convector, heating_infrared_panels, heating_floor_heating
FROM properties
ORDER BY created_at DESC
`

func (q *Queries) ListProperties(ctx context.Context) ([]Property, error) {
	rows, err := q.db.QueryContext(ctx, listProperties)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Humanreadableid,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Floor,
			&i.UserID,
			pq.Array(&i.Images),
			&i.Thumbnail,
			&i.IsFeatured,
			&i.EnergyClass,
			&i.EnergyConsumptionPrimary,
			&i.EnergyEmissionsIndex,
			&i.EnergyConsumptionGreen,
			&i.DestinationResidential,
			&i.DestinationCommercial,
			&i.DestinationOffice,
			&i.DestinationHoliday,
			&i.OtherUtilitiesTerrance,
			&i.OtherUtilitiesServiceToilet,
			&i.OtherUtilitiesUndergroundStorage,
			&i.OtherUtilitiesStorage,
			&i.PropertyTransaction,
			&i.FurnishedNot,
			&i.FurnishedPartially,
			&i.FurnishedComplete,
			&i.FurnishedLuxury,
			&i.InteriorNeedsRenovation,
			&i.InteriorHasRenovation,
			&i.InteriorGoodState,
			&i.HeatingTermoficare,
			&i.HeatingCentralHeating,
			&i.HeatingBuilding,
			&i.HeatingStove,
			&i.HeatingRadiator,
			&i.HeatingOtherElectrical,
			&i.HeatingGasConvector,
			&i.HeatingInfraredPanels,
			&i.HeatingFloorHeating,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}