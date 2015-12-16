package main

import (
	"flag"
	"fmt"
	"github.com/qualiapps/subject/resources"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	ConfigPort *string
	CoapPort   *string
)

func init() {
	ConfigPort = flag.String("port", "8888", "Config server port")
	CoapPort = flag.String("lport", "56083", "CoAP notifier port")

	flag.Parse()
}

func main() {
	var (
		register   = make(chan resources.Resource) // register resource
		deregister = make(chan string)             // remove resource
		event      = make(chan resources.Resource) // incoming event
		exit       = make(chan os.Signal, 1)       // terminate
	)

	coapServ := NewCoapServer(CoapPort)

	signal.Notify(exit,
		os.Interrupt,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go ServConfig(register, deregister, event)
	go coapServ.Start()

	fmt.Printf("CoAP server was started ... OK\n")
	fmt.Printf("Config server :%s\n", *ConfigPort)
	fmt.Printf("CoAP server :%s\n", *CoapPort)

	for {
		select {
		// register resource
		case res := <-register:
			log.Printf("OK.........Resource %s was added.\n", res.Name)
		// remove resource
		case name := <-deregister:
			log.Printf("OK.........Resource %s was deleted.\n", name)
			coapServ.DeregisterResource(name)
		// change resource
		case resource := <-event:
			log.Printf("OK.........Incoming Event %s\n", resource.Name)
			coapServ.Event(resource.Name, resource.Payload)
		// terminate app
		case <-exit:
			go func() {
				coapServ.DeregisterResources()
				log.Printf("OK.........Terminated")
				os.Exit(0)
			}()
		}
	}
}
