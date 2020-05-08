package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math/rand"
	"regexp"
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

var cachedResult = `
GOROOT=/usr/lib/go-1.14 #gosetup
GOPATH=/home/knox/go:/home/knox/Desktop/Programme/GO #gosetup
/usr/lib/go-1.14/bin/go test -c -o /tmp/___BenchmarkTranslateFromNaviCached_in_fwew_lib fwew_lib #gosetup
/tmp/___BenchmarkTranslateFromNaviCached_in_fwew_lib -test.v -test.bench ^BenchmarkTranslateFromNaviCached$ -test.run ^$ #gosetup
goos: linux
goarch: amd64
pkg: fwew_lib
BenchmarkTranslateFromNaviCached
BenchmarkTranslateFromNaviCached/molte
BenchmarkTranslateFromNaviCached/molte-8         	      19	  63364127 ns/op
BenchmarkTranslateFromNaviCached/pepfil
BenchmarkTranslateFromNaviCached/pepfil-8        	      18	  66669865 ns/op
BenchmarkTranslateFromNaviCached/empty
BenchmarkTranslateFromNaviCached/empty-8         	 1219765	       965 ns/op
BenchmarkTranslateFromNaviCached/säpeykiyevatsi
BenchmarkTranslateFromNaviCached/säpeykiyevatsi-8          	       8	 141604206 ns/op
BenchmarkTranslateFromNaviCached/tseng
BenchmarkTranslateFromNaviCached/tseng-8                   	     129	   9330786 ns/op
BenchmarkTranslateFromNaviCached/luyu
BenchmarkTranslateFromNaviCached/luyu-8                    	      58	  22013140 ns/op
BenchmarkTranslateFromNaviCached/seiyi
BenchmarkTranslateFromNaviCached/seiyi-8                   	      34	  37105955 ns/op
BenchmarkTranslateFromNaviCached/zenuyeke
BenchmarkTranslateFromNaviCached/zenuyeke-8                	      21	  55013742 ns/op
BenchmarkTranslateFromNaviCached/verìn
BenchmarkTranslateFromNaviCached/verìn-8                   	      14	  78093592 ns/op
BenchmarkTranslateFromNaviCached/ketsuktswa'
BenchmarkTranslateFromNaviCached/ketsuktswa'-8             	      10	 115187909 ns/op
BenchmarkTranslateFromNaviCached/tìtusaron
BenchmarkTranslateFromNaviCached/tìtusaron-8               	       7	 148665566 ns/op
BenchmarkTranslateFromNaviCached/fayioang
BenchmarkTranslateFromNaviCached/fayioang-8                	      10	 105317353 ns/op
BenchmarkTranslateFromNaviCached/tsasoaiä
BenchmarkTranslateFromNaviCached/tsasoaiä-8                	       7	 144868661 ns/op
BenchmarkTranslateFromNaviCached/tseyä
BenchmarkTranslateFromNaviCached/tseyä-8                   	       5	 230582122 ns/op
BenchmarkTranslateFromNaviCached/oey
BenchmarkTranslateFromNaviCached/oey-8                     	      57	  19856498 ns/op
BenchmarkTranslateFromNaviCached/ngey
BenchmarkTranslateFromNaviCached/ngey-8                    	     150	   7865232 ns/op
BenchmarkTranslateFromNaviCached/tì'usemä
BenchmarkTranslateFromNaviCached/tì'usemä-8                	       9	 114951975 ns/op
BenchmarkTranslateFromNaviCached/wemtswo
BenchmarkTranslateFromNaviCached/wemtswo-8                 	      20	  58203820 ns/op
BenchmarkTranslateFromNaviCached/pawnengsì
BenchmarkTranslateFromNaviCached/pawnengsì-8               	       9	 132390112 ns/op
BenchmarkTranslateFromNaviCached/tsuknumesì
BenchmarkTranslateFromNaviCached/tsuknumesì-8              	       8	 128860672 ns/op
BenchmarkTranslateFromNaviCached/tsamungwrr
BenchmarkTranslateFromNaviCached/tsamungwrr-8              	      12	  88040828 ns/op
BenchmarkTranslateFromNaviCached/tsamsiyu
BenchmarkTranslateFromNaviCached/tsamsiyu-8                	      10	 102504860 ns/op
BenchmarkTranslateFromNaviCached/'ueyä
BenchmarkTranslateFromNaviCached/'ueyä-8                   	       5	 235393020 ns/op
BenchmarkTranslateFromNaviCached/awngeyä
BenchmarkTranslateFromNaviCached/awngeyä-8                 	       5	 236303245 ns/op
BenchmarkTranslateFromNaviCached/fpi
BenchmarkTranslateFromNaviCached/fpi-8                     	     480	   2552541 ns/op
BenchmarkTranslateFromNaviCached/pe
BenchmarkTranslateFromNaviCached/pe-8                      	     590	   2064867 ns/op
PASS

Process finished with exit code 0
`

var UNcachedResult = `
GOROOT=/usr/lib/go-1.14 #gosetup
GOPATH=/home/knox/go:/home/knox/Desktop/Programme/GO #gosetup
/usr/lib/go-1.14/bin/go test -c -o /tmp/___BenchmarkTranslateFromNavi_in_fwew_lib fwew_lib #gosetup
/tmp/___BenchmarkTranslateFromNavi_in_fwew_lib -test.v -test.bench ^BenchmarkTranslateFromNavi$ -test.run ^$ #gosetup
goos: linux
goarch: amd64
pkg: fwew_lib
BenchmarkTranslateFromNavi
BenchmarkTranslateFromNavi/molte
BenchmarkTranslateFromNavi/molte-8         	      16	  70178235 ns/op
BenchmarkTranslateFromNavi/pepfil
BenchmarkTranslateFromNavi/pepfil-8        	      14	  75358431 ns/op
BenchmarkTranslateFromNavi/empty
BenchmarkTranslateFromNavi/empty-8         	 1235800	       949 ns/op
BenchmarkTranslateFromNavi/säpeykiyevatsi
BenchmarkTranslateFromNavi/säpeykiyevatsi-8          	       7	 147990094 ns/op
BenchmarkTranslateFromNavi/tseng
BenchmarkTranslateFromNavi/tseng-8                   	      55	  21025280 ns/op
BenchmarkTranslateFromNavi/luyu
BenchmarkTranslateFromNavi/luyu-8                    	      34	  33348519 ns/op
BenchmarkTranslateFromNavi/seiyi
BenchmarkTranslateFromNavi/seiyi-8                   	      21	  47801298 ns/op
BenchmarkTranslateFromNavi/zenuyeke
BenchmarkTranslateFromNavi/zenuyeke-8                	      19	  66589827 ns/op
BenchmarkTranslateFromNavi/verìn
BenchmarkTranslateFromNavi/verìn-8                   	      13	  87968462 ns/op
BenchmarkTranslateFromNavi/ketsuktswa'
BenchmarkTranslateFromNavi/ketsuktswa'-8             	       9	 120824317 ns/op
BenchmarkTranslateFromNavi/tìtusaron
BenchmarkTranslateFromNavi/tìtusaron-8               	       7	 153693681 ns/op
BenchmarkTranslateFromNavi/fayioang
BenchmarkTranslateFromNavi/fayioang-8                	       9	 115398211 ns/op
BenchmarkTranslateFromNavi/tsasoaiä
BenchmarkTranslateFromNavi/tsasoaiä-8                	       7	 147780615 ns/op
BenchmarkTranslateFromNavi/tseyä
BenchmarkTranslateFromNavi/tseyä-8                   	       5	 235455531 ns/op
BenchmarkTranslateFromNavi/oey
BenchmarkTranslateFromNavi/oey-8                     	      38	  30625738 ns/op
BenchmarkTranslateFromNavi/ngey
BenchmarkTranslateFromNavi/ngey-8                    	      68	  19975270 ns/op
BenchmarkTranslateFromNavi/tì'usemä
BenchmarkTranslateFromNavi/tì'usemä-8                	       9	 119954164 ns/op
BenchmarkTranslateFromNavi/wemtswo
BenchmarkTranslateFromNavi/wemtswo-8                 	      18	  69015712 ns/op
BenchmarkTranslateFromNavi/pawnengsì
BenchmarkTranslateFromNavi/pawnengsì-8               	       8	 135757527 ns/op
BenchmarkTranslateFromNavi/tsuknumesì
BenchmarkTranslateFromNavi/tsuknumesì-8              	       8	 126856501 ns/op
BenchmarkTranslateFromNavi/tsamungwrr
BenchmarkTranslateFromNavi/tsamungwrr-8              	      12	  95917334 ns/op
BenchmarkTranslateFromNavi/tsamsiyu
BenchmarkTranslateFromNavi/tsamsiyu-8                	      10	 110620711 ns/op
BenchmarkTranslateFromNavi/'ueyä
BenchmarkTranslateFromNavi/'ueyä-8                   	       5	 231322433 ns/op
BenchmarkTranslateFromNavi/awngeyä
BenchmarkTranslateFromNavi/awngeyä-8                 	       5	 228520213 ns/op
BenchmarkTranslateFromNavi/fpi
BenchmarkTranslateFromNavi/fpi-8                     	      84	  14591949 ns/op
BenchmarkTranslateFromNavi/pe
BenchmarkTranslateFromNavi/pe-8                      	      78	  14571685 ns/op
PASS

Process finished with exit code 0
`

func randomLangCode() string {
	num := rand.Intn(6)
	return langs[num]
}

func main() {
	var groupCached plotter.Values
	var groupUncached plotter.Values
	var xNames []string

	var validID = regexp.MustCompile(`.+?/(.+?) *\t *(\d+)\t *(\d+) ns/op`)
	//var validID = regexp.MustCompile(`"| `)

	cachedData := validID.FindAllStringSubmatch(cachedResult, -1)

	for _, cachedDatum := range cachedData {
		float, err := strconv.ParseFloat(cachedDatum[3], 64)
		if err == nil {
			groupCached = append(groupCached, float/1000/1000)

		}
		xNames = append(xNames, cachedDatum[1])
	}

	uncachedData := validID.FindAllStringSubmatch(UNcachedResult, -1)

	for _, uncachedDatum := range uncachedData {
		float, err := strconv.ParseFloat(uncachedDatum[3], 64)
		if err == nil {
			groupUncached = append(groupUncached, float/1000/1000)
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Performance TranslateFromNavi"
	p.Y.Label.Text = "Duration (ms)"

	w := vg.Points(20)

	barsCached, err := plotter.NewBarChart(groupCached, w)
	if err != nil {
		panic(err)
	}
	barsCached.LineStyle.Width = vg.Length(0)
	barsCached.Color = plotutil.Color(0)
	barsCached.Offset = -w

	barsUncached, err := plotter.NewBarChart(groupUncached, w)
	if err != nil {
		panic(err)
	}
	barsUncached.LineStyle.Width = vg.Length(0)
	barsUncached.Color = plotutil.Color(1)

	p.Add(barsCached, barsUncached)
	p.Legend.Add("Cached", barsCached)
	p.Legend.Add("Uncached", barsUncached)
	p.Legend.Top = true
	p.NominalX(xNames...)

	if err := p.Save(20*vg.Inch, 8*vg.Inch, "misc/Duration_Cached_Uncached.png"); err != nil {
		panic(err)
	}
}
