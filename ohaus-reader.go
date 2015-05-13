package ohaus

import (
	"bufio"
	serial "github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"time"
)

type Scale struct {
	PortName string
}

type Datum struct {
	time   time.Time
	weight float64
	unit   string
	err    error
}

func (scale Scale) Open() (port *serial.Port, err error) {
	c := &serial.Config{Name: scale.PortName, Baud: 9600}
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
	var d Datum
	if err != nil {
		d.err = err
		c <- d
		return
	}
	for {
		time := time.Now()
		v, err := scale.Read(port)
		if err != nil {
			log.Fatal(err)
		}
		value := strings.Split(strings.Trim(v, " "), " ")
		weight, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			port.Close()
			d.err = err
			c <- d
			return
		}

		d.time = time
		d.weight = weight
		d.unit = value[1]

		c <- d
	}
}
