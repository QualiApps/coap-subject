package main

import (
	"github.com/dustin/go-coap"
	"github.com/qualiapps/subject/resources"
	"log"
	"net"
	"os"
	"strings"
)

type Request struct {
	Data     []byte
	FromAddr *net.UDPAddr
}

/**
 * Processing response
 * @param *net.UDPConn l - connection instance
 * @param dbClient - db instance
 * @param Response response - res data
 */
func ProcessRequest(l *net.UDPConn, request Request) {
	// parse to CoAP struct
	rv := coap.Message{}
	err := rv.UnmarshalBinary(request.Data)
	if err == nil {
		path := string(os.PathSeparator) + strings.Join(rv.Path(), string(os.PathSeparator))
		switch path {
		case resources.WellKnown:
			Discovery(l, request.FromAddr, &rv)
		default:
			if rv.IsObservable() {
				if rv.Option(coap.Observe).(uint32) == 0 {
					log.Println("AAAAAAAAAAa")
				}
				//r := resources.GetAll()
			}

		}

		if rv.IsObservable() {
			// Send ACK
			if rv.IsConfirmable() {
			}
		}
	}
}
