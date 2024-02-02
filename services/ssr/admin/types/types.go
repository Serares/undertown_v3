package types

import (
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/a-h/templ"
)

const (
	DeletePath   = "/delete"
	EditPath     = "/edit"
	SubmitPath   = "/submit"
	ListPath     = "/list"
	ImagesPrefix = "/assets"
	LoginPath    = "/login"
)

type BasicIncludes struct {
	Header         templ.Component
	BannerSection  templ.Component
	Preload        templ.Component
	Navbar         templ.Component
	Footer         templ.Component
	Scripts        templ.Component
	DropzoneScript templ.ComponentScript
}

type EditIncludes struct {
	HandleDeleteButton templ.ComponentScript
}

type SelectInputs struct {
	Value       string
	DisplayName string
}

type SubmitProps struct {
	PropertyTypes       []SelectInputs
	PropertyTransaction []SelectInputs
	FormMethod          string
	FormAction          string
	Message             string
	ErrorMessage        string
}

type EditProps struct {
	PropertyTypes       []SelectInputs
	PropertyTransaction []SelectInputs
	PropertyFeatures    PropertyFeatures
	Property            lite.Property
	FormMethod          string
	FormAction          string
	Images              []string // The images comming from db are a string of image paths separated by ;
	ErrorMessage        string
	SuccessMessage      string
}

type LoginProps struct {
	Message      string
	ErrorMessage string
}

type ListingProperty struct {
	Title              string
	Address            string
	TransactionType    string
	Price              int64
	DisplayPrice       string
	Thumbnail          string
	EditPropertyPath   string
	DeletePropertyPath string
	Surface            int64
	ImagesNumber       int64
	CreatedTime        string
}

type ListingProps struct {
	Properties     []ListingProperty
	ErrorMessage   string
	SuccessMessage string
}

var PropertyTypes = []SelectInputs{
	{
		Value:       "APARTMENT",
		DisplayName: "Apartament",
	},
	{
		Value:       "HOUSE",
		DisplayName: "Casa",
	},
	{
		Value:       "STUDIO",
		DisplayName: "Garsoniera",
	},
	{
		Value:       "LAND",
		DisplayName: "Teren",
	},
}

var PropertyTransactions = []SelectInputs{
	{
		Value:       "0",
		DisplayName: "Vanzare",
	},
	{
		Value:       "1",
		DisplayName: "Chirie",
	},
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
}

type DeleteScriptProps struct {
	DeleteUrl string
}
