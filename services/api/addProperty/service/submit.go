package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/psql"
	"github.com/Serares/undertown_v3/services/api/addProperty/types"
	"github.com/Serares/undertown_v3/services/api/addProperty/util"
	"github.com/google/uuid"
)

// TODO handle the property creation
// handle the checking if the user id exists before persisting the property

// TODO make sure all fields are mapped from repository model
// TODO move this struct to the handler

type Submit struct {
	Log                *slog.Logger
	PropertyRepository *repository.Properties
}

func NewSubmitService(log *slog.Logger, pr *repository.Properties) Submit {
	return Submit{
		Log:                log,
		PropertyRepository: pr,
	}
}

func (ss *Submit) ProcessProperty(ctx context.Context, property *types.POSTProperty) (uuid.UUID, string, error) {
	var propertyId = uuid.New()

	humanReadableId := util.HumanReadableId(psql.TransactionType(property.PropertyTransaction))
	if err := ss.PropertyRepository.Add(ctx, psql.AddPropertyParams{
		ID:              propertyId,
		Humanreadableid: humanReadableId,
		Title:           property.Title,
		Floor:           property.Floor,
		Images:          property.Images,
		// UserID: TODO,
		Thumbnail:                        property.Thumbnail,
		IsFeatured:                       property.IsFeatured,
		EnergyClass:                      property.EnergyClass,
		EnergyConsumptionPrimary:         property.EnergyConsumptionPrimary,
		EnergyEmissionsIndex:             property.EnergyEmissionsIndex,
		EnergyConsumptionGreen:           property.EnergyConsumptionGreen,
		DestinationResidential:           property.DestinationResidential,
		DestinationCommercial:            property.DestinationCommercial,
		DestinationOffice:                property.DestinationOffice,
		DestinationHoliday:               property.DestinationHoliday,
		OtherUtilitiesTerrance:           property.OtherUtilitiesTerrance,
		OtherUtilitiesServiceToilet:      property.OtherUtilitiesServiceToilet,
		OtherUtilitiesUndergroundStorage: property.OtherUtilitiesUndergroundStorage,
		OtherUtilitiesStorage:            property.OtherUtilitiesStorage,
		PropertyTransaction:              psql.TransactionType(property.PropertyTransaction),
		FurnishedNot:                     property.FurnishedNot,
		FurnishedPartially:               property.FurnishedPartially,
		FurnishedComplete:                property.FurnishedComplete,
		FurnishedLuxury:                  property.FurnishedLuxury,
		InteriorNeedsRenovation:          property.InteriorNeedsRenovation,
		InteriorHasRenovation:            property.InteriorHasRenovation,
		InteriorGoodState:                property.InteriorGoodState,
		HeatingTermoficare:               property.HeatingTermoficare,
		HeatingCentralHeating:            property.HeatingCentralHeating,
		HeatingBuilding:                  property.HeatingBuilding,
		HeatingStove:                     property.HeatingStove,
		HeatingRadiator:                  property.HeatingRadiator,
		HeatingOtherElectrical:           property.HeatingOtherElectrical,
		HeatingGasConvector:              property.HeatingGasConvector,
		HeatingInfraredPanels:            property.HeatingInfraredPanels,
		HeatingFloorHeating:              property.HeatingFloorHeating,
		CreatedAt:                        time.Now().UTC(),
		UpdatedAt:                        time.Now().UTC(),
	}); err != nil {
		return uuid.UUID{}, "", fmt.Errorf("error trying to persist the order with error: %v", err)
	}
	return propertyId, humanReadableId, nil
}
