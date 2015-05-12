package main

import (
	"fmt"
	serial "github.com/tarm/serial"
)

type Scale struct {
	port io.ReadWriteCloser
}

func (scale Scale) Open() error {
	c := serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: time.Second * 1}
	port, err := serial.OpenPort(&c)
	scale.port = port
	return scale, err
}

func (scale Scale) Read() {
	port.Write("IP\r\n")
	buf := make([]byte, 128)
	n, err = port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])

}

func main() {
	scale = scale.Scale{}
	err = scale.Open()
	if err != nil {
		log.Fatal(err)
	}
	value, err = scale.Read()
	fmt.Println(value)
}
