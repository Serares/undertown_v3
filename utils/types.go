package utils

import "github.com/golang-jwt/jwt/v5"

type TransactionType int

const (
	Sell TransactionType = iota
	Rent
	Default // using this value to call the ToInt method
)

// Will return SELL or RENT strings
func (t TransactionType) String() string {
	return [...]string{"SELL", "RENT"}[t]
}

// Not sure if this is needed
func (t TransactionType) ToInt(stringType string) TransactionType {
	switch stringType {
	case "SELL":
		return Sell
	case "RENT":
		return Rent
	default:
		return Default
	}
}

type JWTClaims struct {
	Email   string `json:"email"`
	UserId  string `json:"userId"`
	Isadmin bool   `json:"isadmin"`
	IsSsr   bool   `json:"isssr"`
	jwt.RegisteredClaims
}

// This is the "property" field that's received as a multipart form data POST/PUT request
type RequestProperty struct {
	Title               string           `json:"title"`
	IsFeatured          bool             `json:"is_featured"`
	Price               int64            `json:"price"`
	PropertyType        string           `json:"property_type"`
	PropertyDescription string           `json:"property_description"`
	PropertyAddress     string           `json:"property_address"`
	PropertyTransaction string           `json:"property_transaction"` // THIS HAS TO BE SELL || RENT or else the db insertion will fail see 002_properties.sql
	PropertySurface     int64            `json:"property_surface"`
	ImageNames          []string         `json:"images"`         // this is going to be the raw images names from the RAWIMAGES_S3 bucket
	DeletedImages       []string         `json:"deleted_images"` // deleted_images should be prefixed by the hrID handled by the adminSRR
	Features            PropertyFeatures `json:"-"`
	// Features            map[string]interface{} `json:"-"` // This was used before to unmarshal all the fields that come from SSR and are not specifically placed
}

// ❗Those are the property features that can be unmarshaled from
// the Features json string stored on the backend
type PropertyFeatures struct {
	// Floor                            int64  `json:"floor"`
	EnergyClass                      string `json:"energy_class"`
	EnergyConsumptionPrimary         string `json:"energy_consumption_primary"`
	EnergyEmissionsIndex             string `json:"energy_emissions_index"`
	EnergyConsumptionGreen           string `json:"energy_consumption_green"`
	DestinationResidential           bool   `json:"destination_residential"`
	DestinationCommercial            bool   `json:"destination_commercial"`
	DestinationOffice                bool   `json:"destination_office"`
	DestinationHoliday               bool   `json:"destination_holiday"`
	OtherUtilitiesTerrance           bool   `json:"other_utilities_terrance"`
	OtherUtilitiesServiceToilet      bool   `json:"other_utilities_service_toilet"`
	OtherUtilitiesUndergroundStorage bool   `json:"other_utilities_underground_storage"`
	OtherUtilitiesStorage            bool   `json:"other_utilities_storage"`
	FurnishedNot                     bool   `json:"furnished_not"`
	FurnishedPartially               bool   `json:"furnished_partially"`
	FurnishedComplete                bool   `json:"furnished_complete"`
	FurnishedLuxury                  bool   `json:"furnished_luxury"`
	InteriorNeedsRenovation          bool   `json:"interior_needs_renovation"`
	InteriorHasRenovation            bool   `json:"interior_has_renovation"`
	InteriorGoodState                bool   `json:"interior_good_state"`
	HeatingTermoficare               bool   `json:"heating_termoficare"`
	HeatingCentralHeating            bool   `json:"heating_central_heating"`
	HeatingBuilding                  bool   `json:"heating_building"`
	HeatingStove                     bool   `json:"heating_stove"`
	HeatingRadiator                  bool   `json:"heating_radiator"`
	HeatingOtherElectrical           bool   `json:"heating_other_electrical"`
	HeatingGasConvector              bool   `json:"heating_gas_convector"`
	HeatingInfraredPanels            bool   `json:"heating_infrared_panels"`
	HeatingFloorHeating              bool   `json:"heating_floor_heating"`
	Lat                              string `json:"latitude"`
	Lng                              string `json:"longitude"`
}

type SQSDeleteImages struct {
	Images []string `json:"images"`
}

type SQSProcessRawImages struct {
	HumanReadableId string   `json:"humanReadableId"` // this is needed because the raw images are stored without the hrID and the processedImages are stored prefixed with the ID
	Images          []string `json:"images"`
}
