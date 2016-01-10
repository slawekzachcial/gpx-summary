package main
import (
	"flag"
	"fmt"
	"github.com/slawekzachcial/gpxsummary"
)


func main() {
	flag.Parse()

	for _, filePath := range flag.Args() {
		info := gpxsummary.Process(filePath)
		fmt.Printf("%s\n", info.Format())
	}
}

