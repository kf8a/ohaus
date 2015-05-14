package main

import (
	"flag"
	"fmt"
	"github.com/kf8a/ohaus"
)

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	c := make(chan ohaus.Datum)
	scale := ohaus.Scale{PortName: "/dev/ttyUSB0"}
	if test {
		go scale.TestReader(c)
	} else {
		go scale.Reader(c)
	}

	for {
		d := <-c
		fmt.Println(d.Time, d.Weight, d.Unit)
	}
}
