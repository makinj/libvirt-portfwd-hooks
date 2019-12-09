package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Hook struct {
	Type string
}

type Config struct {
	Domains map[string]Hook
}

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
		log.Fatal(fmt.Errorf("Usage: %s <domain> <action>", os.Args[0]))
	}

	//Get arguments
	hookdir := filepath.Dir(os.Args[0])

	configfilename := path.Join(hookdir, "hooks.json")
	configfile, err := os.Open(configfilename)
	if err != nil {
		log.Fatal(fmt.Errorf("Error opening config file %s: %s", configfilename, err))
	}

	configcontents, err := ioutil.ReadAll(configfile)
	if err != nil {
		log.Fatal(fmt.Errorf("Error reading config file %s: %s", configfilename, err))
	}

	var config Config
	err = json.Unmarshal(configcontents, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading config file %s: %s", configfilename, err))
	}

	//domain := os.Args[1]
	//action := os.Args[2]

	log.Printf(strings.Join(os.Args, " "))
}
