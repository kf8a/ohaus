package ohaus

import (
	"bufio"
	"encoding/json"
	serial "github.com/tarm/serial"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Scale struct {
	PortName string
}

type Datum struct {
	Time   time.Time `json:"time"`
	Weight float64   `json:"weight"`
	Unit   string    `json:"unit"`
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
		c <- d
		time.Sleep(2 * time.Second)
	}
}

func (scale Scale) Reader(c chan Datum) {
	f, err := os.OpenFile("backup-data.json", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var d Datum
	for {
		port, err := scale.Open()

		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		for {
			current_time := time.Now()
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

			d.Time = current_time
			d.Weight = weight
			d.Unit = value[1]

			c <- d

			text, err := json.Marshal(d)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, err = f.WriteString(string(text)); err != nil {
				log.Println(err)
				continue
			}

			time.Sleep(60 * time.Second)
		}

	}
}
