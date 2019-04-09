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
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Centroids []struct {
		Point   []float64 `json:"point"`
		XLon    float64   `json:"xLon"`
		YLat    float64   `json:"yLat"`
		GeoHash string    `json:"geoHash"`
	} `json:"centroids"`
}

func (a *admarea) toLoop() *s2.Loop {
	//TODO case where Geometry is not polygon, here we assume it is always a polygon
	var pts []s2.Point
	for _, pt := range a.Geometry.Coordinates[0] {
		//in geojson lon is first at lat second
		pts = append(pts, s2.PointFromLatLng(s2.LatLngFromDegrees(pt[1], pt[0])))
	}
	loop := s2.LoopFromPoints(pts)
	// fmt.Println("Area of Loop before Normalize : ", loop.Area())
	loop.Normalize()
	// fmt.Println("Area of Loop after Normalize: ", loop.Area())
	//sometimes normalize do not work
	if loop.Area() > 1 {
		loop.Invert()
		// fmt.Println("Area of Loop after Invert: ", loop.Area())
	}

	return loop
}

func (as *admareas) s2CellIds(maxCells int) s2.CellUnion {
	//1.build []s2.Loop
	//2. use to build polgon
	//3. approximate Region
	var loops []*s2.Loop
	for _, area := range *as {
		loops = append(loops, area.toLoop())
	}
	p := s2.PolygonFromLoops(loops)
	rc := &s2.RegionCoverer{MaxLevel: 25, MaxCells: maxCells}
	r := s2.Region(p)
	covering := rc.Covering(r)
	return covering
}
