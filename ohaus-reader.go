package main

import (
	"bufio"
	"fmt"
	serial "github.com/tarm/serial"
	"log"
	"os"
	"time"
)

type Scale struct {
}

func (scale Scale) Open() (port *serial.Port, err error) {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	port, err = serial.OpenPort(c)
	return
}

func (scale Scale) Read(port *serial.Port) (value string) {
	port.Write([]byte("IP\r\n"))
	scanner := bufio.NewScanner(port)
	scanner.Scan()

	value = scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return
}

func main() {
	scale := Scale{}
	port, err := scale.Open()
	if err != nil {
		log.Fatal(err)
	}
	for {
		value := scale.Read(port)
		fmt.Println(time.Now(), value)
	}
}
