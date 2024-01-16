// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: properties.sql

package lite

import (
	"context"
	"time"
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
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    )
`

type AddPropertyParams struct {
	ID                               string
	Humanreadableid                  string
	CreatedAt                        time.Time
	UpdatedAt                        time.Time
	Title                            string
	Floor                            int64
	UserID                           string
	Images                           string
	Thumbnail                        string
	IsFeatured                       int64
	EnergyClass                      string
	EnergyConsumptionPrimary         string
	EnergyEmissionsIndex             string
	EnergyConsumptionGreen           string
	DestinationResidential           int64
	DestinationCommercial            int64
	DestinationOffice                int64
	DestinationHoliday               int64
	OtherUtilitiesTerrance           int64
	OtherUtilitiesServiceToilet      int64
	OtherUtilitiesUndergroundStorage int64
	OtherUtilitiesStorage            int64
	PropertyTransaction              string
	FurnishedNot                     int64
	FurnishedPartially               int64
	FurnishedComplete                int64
	FurnishedLuxury                  int64
	InteriorNeedsRenovation          int64
	InteriorHasRenovation            int64
	InteriorGoodState                int64
	HeatingTermoficare               int64
	HeatingCentralHeating            int64
	HeatingBuilding                  int64
	HeatingStove                     int64
	HeatingRadiator                  int64
	HeatingOtherElectrical           int64
	HeatingGasConvector              int64
	HeatingInfraredPanels            int64
	HeatingFloorHeating              int64
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
		arg.Images,
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
WHERE humanReadableId = ?
`

func (q *Queries) DeletePropertyByHumanReadableId(ctx context.Context, humanreadableid string) error {
	_, err := q.db.ExecContext(ctx, deletePropertyByHumanReadableId, humanreadableid)
	return err
}

const getByHumanReadableId = `-- name: GetByHumanReadableId :one
SELECT id, humanreadableid, created_at, updated_at, title, floor, user_id, images, thumbnail, is_featured, energy_class, energy_consumption_primary, energy_emissions_index, energy_consumption_green, destination_residential, destination_commercial, destination_office, destination_holiday, other_utilities_terrance, other_utilities_service_toilet, other_utilities_underground_storage, other_utilities_storage, property_transaction, furnished_not, furnished_partially, furnished_complete, furnished_luxury, interior_needs_renovation, interior_has_renovation, interior_good_state, heating_termoficare, heating_central_heating, heating_building, heating_stove, heating_radiator, heating_other_electrical, heating_gas_convector, heating_infrared_panels, heating_floor_heating
FROM properties
WHERE humanReadableId = ?
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
		&i.Images,
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
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetProperty(ctx context.Context, id string) (Property, error) {
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
		&i.Images,
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

const listFeaturedProperties = `-- name: ListFeaturedProperties :many
SELECT id,
    humanReadableId,
    created_at,
    title,
    thumbnail
FROM properties
where is_featured = 1
ORDER BY created_at DESC
`

type ListFeaturedPropertiesRow struct {
	ID              string
	Humanreadableid string
	CreatedAt       time.Time
	Title           string
	Thumbnail       string
}

func (q *Queries) ListFeaturedProperties(ctx context.Context) ([]ListFeaturedPropertiesRow, error) {
	rows, err := q.db.QueryContext(ctx, listFeaturedProperties)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFeaturedPropertiesRow
	for rows.Next() {
		var i ListFeaturedPropertiesRow
		if err := rows.Scan(
			&i.ID,
			&i.Humanreadableid,
			&i.CreatedAt,
			&i.Title,
			&i.Thumbnail,
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
			&i.Images,
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
