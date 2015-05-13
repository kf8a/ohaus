package main

import (
	"fmt"
	"github.com/kf8a/ohaus"
	"log"
	"time"
)

func main() {
	c := make(chan ohaus.Datum)
	for {
		scale := ohaus.Scale{PortName: "/dev/ttyUSB0"}
		go scale.Reader(c)
		for {
			d := <-c
			if d.Err != nil {
				log.Println(d.Err)
				break
			}
			fmt.Println(d.Time, d.Weight, d.Unit)
		}
		time.Sleep(2 * time.Second)
	}
}
