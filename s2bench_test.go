package s2bench

import (
	"encoding/json"
	"github.com/golang/geo/s2"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func toAdmAreas(fileName string) admareas {
	file, _ := ioutil.ReadFile(fileName)
	var areas admareas
	json.Unmarshal([]byte(file), &areas)
	return areas
}

var idf50, idf100, idf500, idf1000, idf10000, idf100000 s2.CellUnion

const (
	latInside  float64 = 48.85805170891599
	lonInside  float64 = 2.3258399963378906
	latOutside float64 = 48.94911182153499
	lonOutside float64 = 2.4176788330078125
)

func setup() {
	admIDF := toAdmAreas("area_idf_fromadm-area.json")
	idf50 = admIDF.s2CellIds(50)
	idf100 = admIDF.s2CellIds(100)
	idf500 = admIDF.s2CellIds(500)
	idf1000 = admIDF.s2CellIds(1000)
	idf10000 = admIDF.s2CellIds(10000)
	idf100000 = admIDF.s2CellIds(100000)
}
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}
func TestNormalize(t *testing.T) {
	admareas := toAdmAreas("admarea_withproblem.json")
	if len(admareas) != 1 {
		t.Errorf("Expected 1 area in file with problem, got %v\n", len(admareas))
	}
	area := admareas[0].toLoop().Area()
	if area > 1 {
		t.Errorf("Area of %v is too big : %v, expected less than 1", admareas[0].IDCode, area)
	}
}

func TestIdf(t *testing.T) {
	admareas := toAdmAreas("area_idf_fromadm-area.json")
	// admareas := toAdmAreas("admarea_alone.json")
	// admareas := toAdmAreas("admarea_several.json")

	//test below to ensure that the dataset has not been changed
	expectedNbAreas := 23
	if len(admareas) != expectedNbAreas {
		t.Errorf("Expected %v, got %v\n", expectedNbAreas, len(admareas))
	}
	cellIds := admareas.s2CellIds(100)
	var cells []string
	for _, cellID := range cellIds {
		cells = append(cells, cellID.ToToken())
	}

	cellsString := strings.Join(cells, ",")
	// fmt.Println("Idf cells tokens : ", cellsString)
	expectedCellsString := "47e60d2c,47e60d34,47e60d3c,47e60d44,47e60d4c,47e6630d,47e66314,47e6633,47e6635,47e6637,47e66384,47e6639d,47e6643c,47e6645,47e6647,47e664c,47e6653,47e6655,47e66564,47e6656c,47e6657c,47e66584,47e66c2c,47e66c34,47e66c4b,47e66d74,47e66d7c,47e66d9,47e66da4,47e66db4,47e66dbc,47e66dd,47e66df,47e66e4,47e66e84,47e66e8c,47e66ef4,47e66efc,47e66f01,47e66f07,47e66f084,47e66f74,47e66f7c,47e66fc,47e6701,47e6703,47e6705,47e67067,47e6706c,47e67074,47e67175,47e671c,47e6721,47e6723,47e67244,47e672484,47e6725c,47e6727,47e67284,47e6728c,47e67294,47e672b,47e672d,47e672f,47e67a94,47e67a9c,47e67ab,47e67ad,47e67aeaac,47e67b2c,47e67b4d,47e67b54,47e67cac,47e67cb4,47e67cbc,47e67cc4,47e67ccc,47e67cd1ffc,47e67cd3"
	if cellsString != expectedCellsString {
		t.Errorf("IDF Cells doest not match, expected %v, got %v\n", expectedCellsString, cellsString)
	}

}

func callBidderNaive(idf s2.CellUnion, b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !containsBidderImpl(idf, latInside, lonInside) {
			b.Errorf("%v, %v should have been inside IDF, but found outside in naive impl\n", latInside, lonInside)
		}
		if containsBidderImpl(idf, latOutside, lonOutside) {
			b.Errorf("%v, %v should have been outside IDF, but found inside in naive impl\n", latOutside, lonOutside)
		}
	}
}

func callS2Contains(idf s2.CellUnion, b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !containsS2Impl(idf, latInside, lonInside) {
			b.Errorf("%v, %v should have been inside IDF, but found outside in naive impl\n", latInside, lonInside)
		}
		if containsS2Impl(idf, latOutside, lonOutside) {
			b.Errorf("%v, %v should have been outside IDF, but found inside in naive impl\n", latOutside, lonOutside)
		}
	}
}
func BenchmarkBidder50(b *testing.B) {
	callBidderNaive(idf50, b)
}

func BenchmarkBidder100(b *testing.B) {
	callBidderNaive(idf100, b)
}
func BenchmarkBidder500(b *testing.B) {
	callBidderNaive(idf500, b)
}
func BenchmarkBidder1000(b *testing.B) {
	callBidderNaive(idf1000, b)
}
func BenchmarkBidder10000(b *testing.B) {
	callBidderNaive(idf10000, b)
}
func BenchmarkBidder100000(b *testing.B) {
	callBidderNaive(idf100000, b)
}

func BenchmarkS2Contains50(b *testing.B) {
	callS2Contains(idf50, b)
}
func BenchmarkS2Contains100(b *testing.B) {
	callS2Contains(idf100, b)
}
func BenchmarkS2Contains500(b *testing.B) {
	callS2Contains(idf500, b)
}
func BenchmarkS2Contains1000(b *testing.B) {
	callS2Contains(idf1000, b)
}
func BenchmarkS2Contains10000(b *testing.B) {
	callS2Contains(idf10000, b)
}
func BenchmarkS2Contains100000(b *testing.B) {
	callS2Contains(idf100000, b)
}

func Test75001(t *testing.T) {
	//geoJson below long before lat
	area75001 := [][]float64{
		{2.320782077542493, 48.8630784802801},
		{2.325754396571458, 48.86954640751676},
		{2.327877012030475, 48.86986349270285},
		{2.350833552487191, 48.86334427949446},
		{2.350088626519611, 48.861955602681036},
		{2.344559184230591, 48.85399263293169},
		{2.332852142283249, 48.85930664755518},
		{2.320782077542493, 48.8630784802801},
	}
	loop75001 := toLoop(area75001)
	loop75001.Normalize()
	// fmt.Printf("75001 area = %.7f\n", loop75001.Area())
	if loop75001.IsHole() {
		t.Error("loop75001 is not expected to be a hole")
	}
	rc := &s2.RegionCoverer{MaxLevel: 25, MaxCells: 100}
	r := s2.Region(loop75001)
	covering := rc.Covering(r)

	var cells []string
	for _, cellID := range covering {
		cells = append(cells, cellID.ToToken())
	}

	// fmt.Println(strings.Join(cells, ","))

}
