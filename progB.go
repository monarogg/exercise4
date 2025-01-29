package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

		last_message_str := string(buffer[:n])
		last_message_str = strings.TrimSpace(last_message_str)

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

	//buffer := make([]byte, 1024)

	fmt.Println("VALUE: ", value)
	message := value

	for {
		for i := 0; i < 100; i++ {
			message++

			meld := strconv.Itoa(message)
			_, err := conn.Write([]byte(meld))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Sendt: ", meld)
			time.Sleep(1 * time.Second)
		}
	}

}

func drep() {

	time.Sleep(10 * time.Second)
	fmt.Println("Drep nå!")
	os.Exit(0)

}
