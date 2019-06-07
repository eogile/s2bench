package s2bench

import "github.com/golang/geo/s2"

//GeoJSONs is type mapping array of geoJson
type GeoJSONs []GeoJSON

//GeoJSON is type mapping geoJson
type GeoJSON struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func (g *GeoJSON) toLoop() *s2.Loop {
	//TODO case where Geometry is not polygon, here we assume it is always a polygon
	var pts []s2.Point
	for _, pt := range g.Coordinates[0] {
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

func (geoJsons *GeoJSONs) toS2CellIds(maxCells int) s2.CellUnion {
	//1.build []s2.Loop
	//2. use to build polgon
	//3. approximate Region
	var loops []*s2.Loop
	for _, geoJSON := range *geoJsons {
		loops = append(loops, geoJSON.toLoop())
	}
	p := s2.PolygonFromLoops(loops)
	rc := &s2.RegionCoverer{MaxLevel: 25, MaxCells: maxCells}
	r := s2.Region(p)
	covering := rc.Covering(r)
	return covering
}
