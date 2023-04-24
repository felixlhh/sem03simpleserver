package main

import (
	//"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/felixlhh/is105sem03/mycrypt"
)

func celsiusToFahrenheit(temp float64) float64 {

	return (temp * 9 / 5) + 32

}

func konverterSetning(input string) string {

	parts := strings.Split(input, ";")

	if len(parts) != 4 {
		log.Println("Please input valid data")
	}

	tempvalue, err := strconv.ParseFloat(parts[3], 64)
	if err != nil || parts[3] == "" {
		log.Println("Error parsing temperature:", err)
	}

	fahrenheit := celsiusToFahrenheit(tempvalue)
	fahrenheitString := strconv.FormatFloat(fahrenheit, 'f', 1, 64)

	parts[3] = fahrenheitString

	konvertertSetning := strings.Join(parts, ";")

	return konvertertSetning

}

func main() {

	log.Println("test")

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:1234")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					log.Println("Buffer: ", string(buf[:n]))
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					msg := string(dekryptertMelding)
					log.Println("msg :", msg)
					splitmsg := strings.Split(msg, ";")
					log.Println("Splitmsg: ", splitmsg)

					switch splitmsg[0] {
					case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))
					case "Kjevik":
						convertTemp := konverterSetning(msg)
						log.Println(convertTemp)
						kryptertMelding := mycrypt.Krypter([]rune(convertTemp), mycrypt.ALF_SEM03, 4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))
					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}

