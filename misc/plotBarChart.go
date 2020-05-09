package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io/ioutil"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
)

var langs = []string{
	"en",
	"de",
	"et",
	"fr",
	"hu",
	"nl",
	"pl",
	"ru",
	"sv",
}

func randomLangCode() string {
	num := rand.Intn(6)
	return langs[num]
}

func main() {
	type dataStruct struct {
		name     string
		cached   float64
		uncached float64
	}

	var data = map[string]*dataStruct{}
	var err error

	bigBenchmark, err := ioutil.ReadFile("misc/BigBenchmark.txt")
	if err != nil {
		panic(err)
	}

	UNcachedResult := string(bigBenchmark)

	bigBenchmarkCached, err := ioutil.ReadFile("misc/BigBenchmarkCached.txt")
	if err != nil {
		panic(err)
	}

	cachedResult := string(bigBenchmarkCached)

	var validID = regexp.MustCompile(`.+?/(.+?)-\d *\t *(\d+)\t *(\d+) ns/op`)

	cachedData := validID.FindAllStringSubmatch(cachedResult, -1)

	for _, cachedDatum := range cachedData {
		float, err := strconv.ParseFloat(cachedDatum[3], 64)
		if err != nil {
			panic(err)
		}
		dataField := data[cachedDatum[1]]
		if dataField == nil {
			dataField = &dataStruct{}
		}
		dataField.cached = float / 1000 / 1000
		dataField.name = cachedDatum[1]
		data[cachedDatum[1]] = dataField
	}

	uncachedData := validID.FindAllStringSubmatch(UNcachedResult, -1)

	for _, uncachedDatum := range uncachedData {
		float, err := strconv.ParseFloat(uncachedDatum[3], 64)
		if err != nil {
			panic(err)
		}
		dataField := data[uncachedDatum[1]]
		if dataField == nil {
			dataField = &dataStruct{}
		}
		dataField.uncached = float / 1000 / 1000
		dataField.name = uncachedDatum[1]
		data[uncachedDatum[1]] = dataField
	}

	// move data to array instead of struct (arrays can be sorted, maps not)
	var dataArray []*dataStruct
	for _, d := range data {
		dataArray = append(dataArray, d)
	}

	// Sort by slowest uncached
	sort.Slice(dataArray, func(i, j int) bool {
		return math.Abs(dataArray[i].uncached-dataArray[i].cached) > math.Abs(dataArray[j].uncached-dataArray[j].cached)
	})

	dataArray = dataArray[5:]

	var groupCached plotter.Values
	var groupUncached plotter.Values
	var xNames []string

	for i, d := range dataArray {
		groupCached = append(groupCached, d.cached)
		groupUncached = append(groupUncached, d.uncached)
		xNames = append(xNames, d.name)
		if i >= 20 {
			break
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Performance TranslateFromNavi (slowest Uncached)"
	p.Y.Label.Text = "Duration (ms)"

	w := vg.Points(20)

	barsCached, err := plotter.NewBarChart(groupCached, w)
	if err != nil {
		panic(err)
	}
	barsCached.LineStyle.Width = vg.Length(0)
	barsCached.Color = plotutil.Color(1)
	barsCached.Offset = -w

	barsUncached, err := plotter.NewBarChart(groupUncached, w)
	if err != nil {
		panic(err)
	}
	barsUncached.LineStyle.Width = vg.Length(0)
	barsUncached.Color = plotutil.Color(0)

	p.Add(barsCached, barsUncached)
	p.Legend.Add("Cached", barsCached)
	p.Legend.Add("Uncached", barsUncached)
	p.Legend.Top = true
	p.NominalX(xNames...)

	if err := p.Save(20*vg.Inch, 10*vg.Inch, "misc/Duration_Cached_Uncached.png"); err != nil {
		panic(err)
	}

	// Second plot as histogram to show overall difference
	var differenceData plotter.Values

	for _, d := range dataArray {
		differenceData = append(differenceData, d.uncached-d.cached)
	}

	// create second plot
	p2, err := plot.New()
	if err != nil {
		panic(err)
	}
	p2.Title.Text = "Overall Difference"

	// create histogram with our difference7
	// second param is the amount of bars.
	h2, err := plotter.NewHist(differenceData, 32)
	if err != nil {
		panic(err)
	}

	// no normalize, it will just adjust the numbers to under 1 :)
	//h2.Normalize(1)
	p2.Add(h2)

	if err := p2.Save(8*vg.Inch, 4*vg.Inch, "misc/Duration_Overall.png"); err != nil {
		panic(err)
	}
}
