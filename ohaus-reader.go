package main

import (
	"bufio"
	"fmt"
	serial "github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"time"
)

type Scale struct {
}

type Datum struct {
	time   time.Time
	weight float64
	unit   string
}

func (scale Scale) Open() (port *serial.Port, err error) {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	port, err = serial.OpenPort(c)
	return
}

func (scale Scale) Read(port *serial.Port) (value string, err error) {
	port.Write([]byte("IP\r\n"))
	scanner := bufio.NewScanner(port)
	scanner.Scan()
	value = scanner.Text()
	err = scanner.Err()
	return
}

func (scale Scale) Reader(c chan Datum) {
	port, err := scale.Open()
	if err != nil {
		log.Fatal(err)
	}
	for {
		v, err := scale.Read(port)
		if err != nil {
			log.Fatal(err)
		}
		value := strings.Split(strings.Trim(v, " "), " ")
		weight, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		d := Datum{
			time:   time.Now(),
			weight: weight,
			unit:   value[1],
		}

		// fmt.Println(d)
		c <- d
	}
}

func main() {
	c := make(chan Datum)
	scale := Scale{}
	go scale.Reader(c)
	for {
		d := <-c
		fmt.Println(d.time, d.weight, d.unit)
	}
}
