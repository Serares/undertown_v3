package types

import "github.com/Serares/undertown_v3/repositories/repository/types"

const ErrorMethodNotSupported string = "http method not supported"

type POSTSuccessResponse struct {
	PropertyId      string
	HumanReadableId string
}

// This is the "property" field that's send as a multipart form data POST request
type RequestProperty struct {
	Title               string                `json:"title"`
	IsFeatured          bool                  `json:"is_featured"`
	Price               int                   `json:"price"`
	PropertyType        string                `json:"property_type"`
	PropertyDescription string                `json:"property_description"`
	PropertyAddress     string                `json:"property_address"`
	PropertyTransaction types.TransactionType `json:"property_transaction"`
	PropertySurface     int                   `json:"property_surface"`
	// Floor int64  `json:"floor"`
	// Images []string `json:"images"`
	// Thumbnail string   `json:"thumbnail"`
	// UserId                           string   `json:"userid"` // this is added by the authorizer to the request context
	// EnergyClass                      string `json:"energyClass"`
	// EnergyConsumptionPrimary         string `json:"energyConsumptionPrimary"`
	// EnergyEmissionsIndex             string `json:"energyEmissionsIndex"`
	// EnergyConsumptionGreen           string `json:"energyConsumptionGreen"`
	// DestinationResidential           bool   `json:"destinationResidential"`
	// DestinationCommercial            bool   `json:"destinationCommercial"`
	// DestinationOffice                bool   `json:"destinationOffice"`
	// DestinationHoliday               bool   `json:"destinationHoliday"`
	// OtherUtilitiesTerrance           bool   `json:"otherUtilitiesTerrance"`
	// OtherUtilitiesServiceToilet      bool   `json:"otherUtilitiesServiceToilet"`
	// OtherUtilitiesUndergroundStorage bool   `json:"otherUtilitiesUndergroudStorage"`
	// OtherUtilitiesStorage            bool   `json:"otherUtilitiesStorage"`
	// // TODO this might be a issue when trying to unmarshal
	// FurnishedNot            bool `json:"furnishedNot"`
	// FurnishedPartially      bool `json:"furnishedPartially"`
	// FurnishedComplete       bool `json:"furnishedComplete"`
	// FurnishedLuxury         bool `json:"furnishedLuxury"`
	// InteriorNeedsRenovation bool `json:"fnteriorNeedsRenovation"`
	// InteriorHasRenovation   bool `json:"fnteriorHasRenovation"`
	// InteriorGoodState       bool `json:"fnteriorGoodState"`
	// HeatingTermoficare      bool `json:"featingTermoficare"`
	// HeatingCentralHeating   bool `json:"featingCentralHeating"`
	// HeatingBuilding         bool `json:"featingBuilding"`
	// HeatingStove            bool `json:"heatingStove"`
	// HeatingRadiator         bool `json:"heatingRadiator"`
	// HeatingOtherElectrical  bool `json:"heatingOtherElectrical"`
	// HeatingGasConvector     bool `json:"heatingGasConvector"`
	// HeatingInfraredPanels   bool `json:"heatingInfraredPanels"`
	// HeatingFloorHeating     bool `json:"heatingFloorHeating"`
}

type RequestFeatures struct {
	Floor                            int64
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

type POSTFormData struct {
	Title  string
	Images []string
	Floor  int64
}
