package main
import (
	"flag"
	"fmt"
	"github.com/slawekzachcial/gpxsummary"
	"sort"
	"os"
)


func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-t] gpxfile1 ...\n\nPrints GPX tracks summaries\n", os.Args[0])
		flag.PrintDefaults()
	}

	asTable := flag.Bool("t", false, "Show data as table")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	tracks := make([]gpxsummary.TrackInfo, flag.NArg())

	for i, filePath := range flag.Args() {
		tracks[i] = gpxsummary.Process(filePath)
	}

	sort.Sort(gpxsummary.TrackInfoArray(tracks))

	for k, info := range tracks {
		if *asTable && k == 0 {
			fmt.Println(info.FormatHeader())
		}
		if !*asTable && k > 0 {
			fmt.Println("---")
		}

		if *asTable {
			fmt.Println(info.FormatRow())
		} else {
			fmt.Println(info.Format())
		}
	}
}

