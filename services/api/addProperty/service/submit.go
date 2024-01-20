package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
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

func arrayToString(images []string) string {
	return strings.Join(images, ",")
}

func booleanToInt(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func (ss *Submit) ProcessProperty(ctx context.Context, property *types.POSTProperty) (string, string, error) {
	var propertyId = uuid.New().String()

	humanReadableId := util.HumanReadableId(property.PropertyTransaction)
	if err := ss.PropertyRepository.Add(ctx, lite.AddPropertyParams{
		ID: propertyId,
		// UserID: , TODO
		Humanreadableid:                  humanReadableId,
		Title:                            property.Title,
		Floor:                            property.Floor,
		Images:                           arrayToString(property.Images),
		Thumbnail:                        property.Thumbnail,
		IsFeatured:                       booleanToInt(property.IsFeatured),
		EnergyClass:                      property.EnergyClass,
		EnergyConsumptionPrimary:         property.EnergyConsumptionPrimary,
		EnergyEmissionsIndex:             property.EnergyEmissionsIndex,
		EnergyConsumptionGreen:           property.EnergyConsumptionGreen,
		DestinationResidential:           booleanToInt(property.DestinationResidential),
		DestinationCommercial:            booleanToInt(property.DestinationCommercial),
		DestinationOffice:                booleanToInt(property.DestinationOffice),
		DestinationHoliday:               booleanToInt(property.DestinationHoliday),
		OtherUtilitiesTerrance:           booleanToInt(property.OtherUtilitiesTerrance),
		OtherUtilitiesServiceToilet:      booleanToInt(property.OtherUtilitiesServiceToilet),
		OtherUtilitiesUndergroundStorage: booleanToInt(property.OtherUtilitiesUndergroundStorage),
		OtherUtilitiesStorage:            booleanToInt(property.OtherUtilitiesStorage),
		PropertyTransaction:              property.PropertyTransaction.String(),
		PropertyDescription:              property.PropertyDescription,
		PropertyType:                     property.PropertyType,
		PropertyAddress:                  property.PropertyAddress,
		PropertySurface:                  int64(property.PropertySurface),
		Price:                            int64(property.Price),
		FurnishedNot:                     booleanToInt(property.FurnishedNot),
		FurnishedPartially:               booleanToInt(property.FurnishedPartially),
		FurnishedComplete:                booleanToInt(property.FurnishedComplete),
		FurnishedLuxury:                  booleanToInt(property.FurnishedLuxury),
		InteriorNeedsRenovation:          booleanToInt(property.InteriorNeedsRenovation),
		InteriorHasRenovation:            booleanToInt(property.InteriorHasRenovation),
		InteriorGoodState:                booleanToInt(property.InteriorGoodState),
		HeatingTermoficare:               booleanToInt(property.HeatingTermoficare),
		HeatingCentralHeating:            booleanToInt(property.HeatingCentralHeating),
		HeatingBuilding:                  booleanToInt(property.HeatingBuilding),
		HeatingStove:                     booleanToInt(property.HeatingStove),
		HeatingRadiator:                  booleanToInt(property.HeatingRadiator),
		HeatingOtherElectrical:           booleanToInt(property.HeatingOtherElectrical),
		HeatingGasConvector:              booleanToInt(property.HeatingGasConvector),
		HeatingInfraredPanels:            booleanToInt(property.HeatingInfraredPanels),
		HeatingFloorHeating:              booleanToInt(property.HeatingFloorHeating),
		CreatedAt:                        time.Now().UTC(),
		UpdatedAt:                        time.Now().UTC(),
	}); err != nil {
		return "", "", fmt.Errorf("error trying to persist the order with error: %v", err)
	}
	return propertyId, humanReadableId, nil
}
