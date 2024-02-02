package types

const ErrorMethodNotSupported string = "http method not supported"

type POSTSuccessResponse struct {
	PropertyId      string
	HumanReadableId string
}

type RequestFeatures struct {
	// Floor                            int64 This is not present in the UI template
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
