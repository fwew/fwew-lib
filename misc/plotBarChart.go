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

type dataStruct struct {
	name     string
	cached   float64
	uncached float64
}

func main() {
	// TranslateFromNavi test results !!!

	dataArray := readBenchResultFiles("misc/BigBenchmarkCached.txt", "misc/BigBenchmark.txt")

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

	printBar(groupCached, groupUncached, xNames, "misc/Duration_Cached_Uncached.png")

	// Second plot as histogram to show overall difference
	var differenceData plotter.Values

	for _, d := range dataArray {
		differenceData = append(differenceData, d.uncached-d.cached)
	}
	printHistogram(differenceData, "misc/Duration_Overall.png")

	var cachedValues plotter.Values
	for _, d := range dataArray {
		cachedValues = append(cachedValues, d.cached)
	}
	printHistogram(cachedValues, "misc/Performance_Cached.png")

	// TranslateToNavi test results !!!
	dataArrayEnglish := readBenchResultFiles("misc/BigBenchmarkEnglishCached.txt", "misc/BigBenchmarkEnglish.txt")

	// Sort by slowest uncached
	sort.Slice(dataArrayEnglish, func(i, j int) bool {
		return math.Abs(dataArrayEnglish[i].uncached-dataArrayEnglish[i].cached) > math.Abs(dataArrayEnglish[j].uncached-dataArrayEnglish[j].cached)
	})

	var groupCachedEnglish plotter.Values
	var groupUncachedEnglish plotter.Values
	var xNamesEnglish []string

	for i, d := range dataArrayEnglish {
		groupCachedEnglish = append(groupCachedEnglish, d.cached)
		groupUncachedEnglish = append(groupUncachedEnglish, d.uncached)
		xNamesEnglish = append(xNamesEnglish, d.name)
		if i >= 20 {
			break
		}
	}

	printBar(groupCachedEnglish, groupUncachedEnglish, xNamesEnglish, "misc/Duration_Cached_Uncached_English.png")

	// Second plot as histogram to show overall difference
	var differenceDataEnglish plotter.Values

	for _, d := range dataArrayEnglish {
		differenceDataEnglish = append(differenceDataEnglish, d.uncached-d.cached)
	}
	printHistogram(differenceDataEnglish, "misc/Duration_Overall_English.png")

	var cachedValuesEnglish plotter.Values
	for _, d := range dataArrayEnglish {
		cachedValuesEnglish = append(cachedValuesEnglish, d.cached)
	}
	printHistogram(cachedValuesEnglish, "misc/Performance_Cached_English.png")
}

func readBenchResultFiles(cachedFile, uncachedFile string) []*dataStruct {
	var data = map[string]*dataStruct{}

	bigBenchmark, err := ioutil.ReadFile(cachedFile)
	if err != nil {
		panic(err)
	}

	UNcachedResult := string(bigBenchmark)

	bigBenchmarkCached, err := ioutil.ReadFile(uncachedFile)
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
	return dataArray
}

func printHistogram(values plotter.Values, filename string) {
	// create third plot
	p3, err := plot.New()
	if err != nil {
		panic(err)
	}
	p3.Title.Text = "Overall Performance (Cached)"
	p3.X.Label.Text = "ms"

	h3, err := plotter.NewHist(values, 32)
	if err != nil {
		panic(err)
	}

	p3.Add(h3)

	if err := p3.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func printBar(cached, uncached plotter.Values, xNames []string, filename string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Performance (sort by highest diff)"
	p.Y.Label.Text = "Duration (ms)"

	w := vg.Points(20)

	barsCached, err := plotter.NewBarChart(cached, w)
	if err != nil {
		panic(err)
	}
	barsCached.LineStyle.Width = vg.Length(0)
	barsCached.Color = plotutil.Color(1)
	barsCached.Offset = -w

	barsUncached, err := plotter.NewBarChart(uncached, w)
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

	if err := p.Save(20*vg.Inch, 10*vg.Inch, filename); err != nil {
		panic(err)
	}
}
