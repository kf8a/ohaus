package ohaus

import (
	"bufio"
	serial "github.com/tarm/serial"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Scale struct {
	PortName string
}

type Datum struct {
	Time   time.Time
	Weight float64
	Unit   string
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

func (scale Scale) TestReader(c chan Datum) {
	var d Datum
	d.Unit = "kg"
	for {
		d.Time = time.Now()
		d.Weight = rand.Float64()
		time.Sleep(1 * time.Second)
		c <- d
	}
}

func (scale Scale) Reader(c chan Datum) {
	var d Datum
	for {
		port, err := scale.Open()

		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		for {
			time := time.Now()
			v, err := scale.Read(port)
			if err != nil {
				port.Close()
				log.Println(err)
				break
			}
			value := strings.Split(strings.Trim(v, " "), " ")
			weight, err := strconv.ParseFloat(value[0], 64)
			if err != nil || len(value) < 2 {
				port.Close()
				log.Println(err)
				break
			}

			d.Time = time
			d.Weight = weight
			d.Unit = value[1]

			c <- d
		}

	}
}
