package s2bench

import (
	"github.com/golang/geo/s2"
)

//mimic gitlab.com/zerobacon/rtb/lib/go-rtb/rtb/filters.go line 194
func containsBidderImpl(cells s2.CellUnion, lat float64, lon float64) bool {
	cell := s2.CellFromLatLng(s2.LatLngFromDegrees(lat, lon))
	ok := false
	for _, cellID := range cells {
		ok = cellID.Contains(cell.ID())
		if ok {
			return ok
		}
	}
	return ok
}

//implementation using S2 methods without loop
func containsS2Impl(cells s2.CellUnion, lat float64, lon float64) bool {
	pt := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lon))
	return cells.ContainsPoint(pt)
}

func toLoop(points [][]float64) *s2.Loop {
	var pts []s2.Point
	for _, pt := range points {
		pts = append(pts, s2.PointFromLatLng(s2.LatLngFromDegrees(pt[1], pt[0])))
	}
	return s2.LoopFromPoints(pts)
}
