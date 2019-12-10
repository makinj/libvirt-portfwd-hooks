package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/coreos/go-iptables/iptables"
)

type HostPort struct {
	Interface string
	Ip        string
	Ports     string
}

type PortForward struct {
	Protocol    string
	Source      HostPort
	Destination HostPort
}

type Action string

func (portfwd PortForward) HandleEvent(action Action) error {
	log.Println("forwarding %s port %s to %s", portfwd.Protocol, portfwd.Source.Ports, portfwd.Destination.Ip)

	natrulespec := []string{"-p", portfwd.Protocol, "--dport", portfwd.Source.Ports, "-i", portfwd.Source.Interface, "-j", "DNAT", "--to-ports", portfwd.Destination.Ports, "--to", portfwd.Destination.Ip, "-o", portfwd.Destination.Interface}
	filterrulespec := []string{"-p", portfwd.Protocol, "--dport", portfwd.Destination.Ports, "-d", portfwd.Destination.Ip, "-i", portfwd.Source.Interface, "-o", portfwd.Destination.Interface, "-m", "state", "--state", "NEW,ESTABLISHED,RELATED", "-j", "ACCEPT"}

	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	if action == "start" {
		err = ipt.Append("nat", "PREROUTING", natrulespec...)
		if err != nil {
			log.Println(err)
		}
		err = ipt.Insert("filter", "FORWARD", 1, filterrulespec...)
		if err != nil {
			log.Println(err)
		}
	} else if action == "stopped" {
		err = ipt.Delete("nat", "PREROUTING", natrulespec...)
		if err != nil {
			log.Println(err)
		}
		err = ipt.Delete("filter", "FORWARD", filterrulespec...)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

type DomainId string

type Config struct {
	Domains map[DomainId][]PortForward
}

func HandleEvent(domain DomainId, action Action, config Config) error {
	log.Printf("Got action %s for domain %s", action, domain)
	forwardlist, ok := config.Domains[domain]

	if !ok {
		log.Printf("No portfwds registered for '%s'", domain)
		for n := range config.Domains {
			log.Printf("portfwd IS registered for '%s'", n)
		}

		return nil
	}

	for _, portfwd := range forwardlist {
		err := portfwd.HandleEvent(action)
		if err != nil {
			log.Println(err)
		}
	}

	//log.Println(hook.Type)
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
	domain := DomainId(os.Args[1])
	action := Action(os.Args[2])

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
	log.Println(configcontents)

	var config Config
	err = json.Unmarshal(configcontents, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading config file %s: %s", configfilename, err))
	}
	log.Println(config)

	// Handle event
	err = HandleEvent(domain, action, config)
	if err != nil {
		log.Fatal(fmt.Errorf("Error processing action %s for domain %s: %s", action, domain, err))
	}
}
