package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	//Setup logs
	f, err := os.OpenFile("/tmp/libvirt-portfwd-hooks.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Verify arguments
	if len(os.Args) < 3 {
		log.Fatal(fmt.Sprintf("Usage: %s <domain> <action>", os.Args[0]))
	}

	//Get arguments
	//domain := os.Args[1]
	//action := os.Args[2]

	log.Printf(strings.Join(os.Args, " "))
}
