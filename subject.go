package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/qualiapps/subject/resources"
	"log"
	"net"
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
		listener   = make(chan *net.UDPConn)       // entry point, to listen notifications
		register   = make(chan resources.Resource) // register resource
		deregister = make(chan string)             // remove resource
		event      = make(chan resources.Resource) // incoming event
		handler    = make(chan Request)            // incoming event
		exit       = make(chan os.Signal, 1)       // terminate
	)

	connString := HostPort{Net, ":" + *CoapPort}

	signal.Notify(exit,
		os.Interrupt,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go ServConfig(register, deregister, event)
	go ServCoap(listener, handler, connString)

	fmt.Printf("CoAP server was started ... OK\n")
	fmt.Printf("Config server :%s\n", *ConfigPort)
	fmt.Printf("CoAP server :%s\n", *CoapPort)

	l := <-listener

	for {
		select {
		// register resource
		case res := <-register:
			log.Printf("OK.........Resource %s was added.\n", res.Name)
		// remove resource
		case name := <-deregister:
			log.Printf("OK.........Resource %s was deleted.\n", name)
			or := observableList[name]
			SendDeregister(l, name, or)
		// change resource
		case resource := <-event:
			log.Printf("OK.........Incoming Event %s\n", resource.Name)
			or := observableList[resource.Name]
			if or != nil {
				for _, r := range or {
					SendNotification(l, r, coap.Content, resource.Payload)
				}
			}
		// incoming handler
		case request := <-handler:
			go ProcessRequest(l, request)
		// terminate app
		case <-exit:
			go func() {
				for route, or := range observableList {
					SendDeregister(l, route, or)
				}
				log.Printf("OK.........Terminated")
				os.Exit(0)
			}()
		}
	}
}

func SendDeregister(l *net.UDPConn, route string, or []*Observation) {
	if or != nil {
		for _, r := range or {
			SendNotification(l, r, coap.NotFound, "")
		}
		log.Printf("OK.........Observation %s was deleted.\n", route)
		delete(observableList, route)
	}

}
