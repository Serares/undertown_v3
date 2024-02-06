package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
)

func UnmarshalProperty(jsonString []byte) (lite.Property, types.PropertyFeatures, error) {
	var marshalingErrors error
	// do the json unmarshal
	var property utils.RequestProperty // using the type from the addProperty service
	marshalingErrors = json.Unmarshal(jsonString, &property)
	if marshalingErrors != nil {
		return lite.Property{}, types.PropertyFeatures{}, marshalingErrors
	}
	// Have to unmarshal the property features into the types.PropertyFeatures struct
	// ❗but the features are unmarshaled first as a map[string]interface{}
	var features types.PropertyFeatures
	// featuresAsJsonString, err := json.Marshal(&property.Features)
	// if err != nil {
	// 	return lite.Property{}, types.PropertyFeatures{}, err
	// }
	marshalingErrors = json.Unmarshal(jsonString, &features)
	if marshalingErrors != nil {
		return lite.Property{}, types.PropertyFeatures{}, marshalingErrors
	}

	return lite.Property{
		Title:               property.Title,
		Price:               property.Price,
		IsFeatured:          utils.BoolToInt(property.IsFeatured),
		PropertyType:        property.PropertyType,
		PropertyDescription: property.PropertyDescription,
		PropertyAddress:     property.PropertyAddress,
		PropertyTransaction: fmt.Sprintf("%d", property.PropertyTransaction),
		PropertySurface:     int64(property.PropertySurface),
		Features:            "", // features can be empty because it's already been unmarshaled
	}, features, nil
}
