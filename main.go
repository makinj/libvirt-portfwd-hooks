package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Hook struct {
	Type string
}

type DomainName string

type Config struct {
	Domains map[DomainName]Hook
}

func HandleEvent(domain DomainName, action string, config Config) error {
	log.Printf("Got action %s for domain %s", action, domain)
	hook := config.Domains[domain]
	log.Println(hook.Type)
	return nil
}

func main() {

	//Setup logs
	logfile, err := os.OpenFile("/var/log/libvirt-portfwd-hooks.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	//Verify arguments
	if len(os.Args) < 3 {
		log.Fatal(fmt.Errorf("Usage: %s <domain> <action>", os.Args[0]))
	}

	//Get arguments
	hookdir := filepath.Dir(os.Args[0])
	domain := DomainName(os.Args[1])
	action := os.Args[2]

	//load config
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

	// Handle event
	err = HandleEvent(domain, action, config)
	if err != nil {
		log.Fatal(fmt.Errorf("Error processing action %s for domain %s: %s", action, domain, err))
	}
}
