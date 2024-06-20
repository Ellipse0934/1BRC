package data

import (
    _ "embed"
    "encoding/json"
)

// TODO: Consider using go:generate to generate a go file with this information
// injected in the final executable

//go:embed extended.json
var extendedset []byte

//go:embed default.json
var defaultset []byte

type CityList struct {
    Cities []string `json:"Cities"`
}

func GetExtended() []string {
    var X CityList
    json.Unmarshal(extendedset, &X)
    return X.Cities   
}

func GetDefault() []string {
    var X CityList
    json.Unmarshal(defaultset, &X)
    return X.Cities   
}
