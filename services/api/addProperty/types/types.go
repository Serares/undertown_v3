package types

const ErrorMethodNotSupported string = "http method not supported"

type POSTSuccessResponse struct {
	PropertyId      string
	HumanReadableId string
}

type POSTProperty struct {
	Title                            string   `json:"title"`
	Floor                            int32    `json:"floor"`
	Images                           []string `json:"images"`
	Thumbnail                        string   `json:"thumbnail"`
	IsFeatured                       bool     `json:"isFeatured"`
	EnergyClass                      string   `json:"energyClass"`
	EnergyConsumptionPrimary         string   `json:"energyConsumptionPrimary"`
	EnergyEmissionsIndex             string   `json:"energyEmissionsIndex"`
	EnergyConsumptionGreen           string   `json:"energyConsumptionGreen"`
	DestinationResidential           bool     `json:"destinationResidential"`
	DestinationCommercial            bool     `json:"destinationCommercial"`
	DestinationOffice                bool     `json:"destinationOffice"`
	DestinationHoliday               bool     `json:"destinationHoliday"`
	OtherUtilitiesTerrance           bool     `json:"otherUtilitiesTerrance"`
	OtherUtilitiesServiceToilet      bool     `json:"otherUtilitiesServiceToilet"`
	OtherUtilitiesUndergroundStorage bool     `json:"otherUtilitiesUndergroudStorage"`
	OtherUtilitiesStorage            bool     `json:"OtherUtilitiesStorage"`
	// TODO this might be a issue when trying to unmarshal
	PropertyTransaction     string `json:"propertyTransaction"`
	FurnishedNot            bool   `json:"furnishedNot"`
	FurnishedPartially      bool   `json:"furnishedPartially"`
	FurnishedComplete       bool   `json:"furnishedComplete"`
	FurnishedLuxury         bool   `json:"furnishedLuxury"`
	InteriorNeedsRenovation bool   `json:"fnteriorNeedsRenovation"`
	InteriorHasRenovation   bool   `json:"fnteriorHasRenovation"`
	InteriorGoodState       bool   `json:"fnteriorGoodState"`
	HeatingTermoficare      bool   `json:"featingTermoficare"`
	HeatingCentralHeating   bool   `json:"featingCentralHeating"`
	HeatingBuilding         bool   `json:"featingBuilding"`
	HeatingStove            bool   `json:"heatingStove"`
	HeatingRadiator         bool   `json:"heatingRadiator"`
	HeatingOtherElectrical  bool   `json:"heatingOtherElectrical"`
	HeatingGasConvector     bool   `json:"heatingGasConvector"`
	HeatingInfraredPanels   bool   `json:"heatingInfraredPanels"`
	HeatingFloorHeating     bool   `json:"heatingFloorHeating"`
}
