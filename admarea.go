package s2bench

import (
	"github.com/golang/geo/s2"
)

type admareas []admarea

type admarea struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Class                   string `json:"_class"`
	ZoneName                string `json:"zoneName"`
	ZoneCountryCodeLevel0   string `json:"zoneCountryCodeLevel0"`
	ZoneCountryName0        string `json:"zoneCountryName0"`
	ZoneCodeLevel1          string `json:"zoneCodeLevel1"`
	ZoneNameLevel1          string `json:"zoneNameLevel1"`
	ZoneCodeLevel2          string `json:"zoneCodeLevel2"`
	ZoneNameLevel2          string `json:"zoneNameLevel2"`
	ZoneCodeLevel3          string `json:"zoneCodeLevel3"`
	ZoneNameLevel3          string `json:"zoneNameLevel3"`
	AdministrativeAreaType  string `json:"administrativeAreaType"`
	AdministrativeAreaLevel int    `json:"administrativeAreaLevel"`
	IsCountrySharpestLevel  bool   `json:"isCountrySharpestLevel"`
	Centroid                struct {
		Point   []float64 `json:"point"`
		XLon    float64   `json:"xLon"`
		YLat    float64   `json:"yLat"`
		GeoHash string    `json:"geoHash"`
	} `json:"centroid"`
	RadiusMeter            int    `json:"radiusMeter"`
	GoogleFormattedAddress string `json:"googleFormattedAddress"`
	IDCode                 string `json:"idCode"`
	Type                   string `json:"type"`
	Properties             struct {
		Class             string `json:"_class"`
		PopulationDetails []struct {
			Population int `json:"population"`
		} `json:"populationDetails"`
		Name string `json:"name"`
	} `json:"properties"`
	Bbox struct {
		SouthWest struct {
			Point []float64 `json:"point"`
			XLon  float64   `json:"xLon"`
			YLat  float64   `json:"yLat"`
		} `json:"southWest"`
		NorthEast struct {
			Point []float64 `json:"point"`
			XLon  float64   `json:"xLon"`
			YLat  float64   `json:"yLat"`
		} `json:"northEast"`
	} `json:"bbox"`
	Geometry  GeoJSON `json:"geometry"`
	Centroids []struct {
		Point   []float64 `json:"point"`
		XLon    float64   `json:"xLon"`
		YLat    float64   `json:"yLat"`
		GeoHash string    `json:"geoHash"`
	} `json:"centroids"`
}

func (a *admarea) toLoop() *s2.Loop {
	return a.Geometry.toLoop()
}

func (as *admareas) s2CellIds(maxCells int) s2.CellUnion {
	// make geoJsons from all admarea.Geometry
	var geometries GeoJSONs
	for _, area := range *as {
		geometries = append(geometries, area.Geometry)
	}
	return geometries.toS2CellIds(maxCells)
}
