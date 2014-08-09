package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	Socket string
	Port string
}

func main() {
	// { "Socket" : "tcp", "Port" : "1337" }
	configuration, err := ioutil.ReadFile("godoor.conf")
	if err != nil {
		panic(err)
	}

	Conf := new(Configuration)

	err = json.Unmarshal(configuration, Conf)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", Conf);

}
