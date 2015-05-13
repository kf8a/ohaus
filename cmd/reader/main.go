package main

import (
	"fmt"
	"ohausreader"
	"time"
)

func main() {
	c := make(chan Datum)
	for {
		scale := Scale{PortName: "/dev/ttyUSB0"}
		go scale.Reader(c)
		for {
			d := <-c
			if d.err != nil {
				log.Println(d.err)
				break
			}
			fmt.Println(d.time, d.weight, d.unit)
		}
		time.Sleep(2 * time.Second)
	}
}
