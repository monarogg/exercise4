package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const port = ":40000"

func main() {

	go drep()

	value, err := listenUDP()
	if err != nil {
		log.Fatal(err)
	}

	exec.Command("gnome-terminal", "--", "go", "run", "progB.go").Run()


	sendUDP(value)

}

func listenUDP() (int, error) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return 0, fmt.Errorf("Error with resolving adress")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return 0, fmt.Errorf("Error with setting up socket")
	}
	defer conn.Close()

	fmt.Println("Listening to broadcast on UDP port")

	err = conn.SetReadDeadline(time.Now().Add(4 * time.Second))
	if err != nil {
		return 0, fmt.Errorf("Error with setting timout")
	}

	buffer := make([]byte, 1024)
	var last_message int

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("No data received within time limit", err)
			break
		}

		// muligens unødvendig:
		if n == 0 {
			fmt.Println("UDP-package empty")
			continue
		}
		last_message_str := string(buffer[:n])

		last_message, err = strconv.Atoi(last_message_str)
		if err != nil {
			 	return 0, fmt.Errorf("Error with converting message")
			}

		fmt.Println("last_message: ", last_message)

	}
	return last_message, nil

}

func sendUDP(value int) {
	broadcastAdrr, err := net.ResolveUDPAddr("udp", "255.255.255.255"+port)
	if err != nil {
		fmt.Println("Error with resolving adress")
		return
	}

	conn, err := net.DialUDP("udp", nil, broadcastAdrr)
	if err != nil {
		fmt.Println("Error with creating UDP-connection")
		return
	}

	defer conn.Close()

	buffer := make([]byte, 1024)

	// message, err := strconv.Atoi(value)
	// if err != nil {
	// 	fmt.Println("Error with converting message")
	// 	return
	// }


	fmt.Println("MESSAGE: ", value)
	message := value

	for {
		for i := 1; i < 100; i++ {
			message += i

			meld := strconv.Itoa(i)
			_, err := conn.Write(buffer)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Sendt: ", meld)
			time.Sleep(1 * time.Second)
		}
		//os.Exit(0)
	}

}

func drep() {

	time.Sleep(10 * time.Second)
	fmt.Println("Drep nå!")
	os.Exit(0)

}
