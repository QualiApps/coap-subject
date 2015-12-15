package main

import (
	"github.com/dustin/go-coap"
	"github.com/qualiapps/subject/resources"
	"github.com/qualiapps/subject/utils"
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
		route, observable := checkRoute(path)
		switch path {
		case resources.WellKnown:
			Discovery(l, request.FromAddr, &rv)
		case route:
			if rv.IsObservable() {
				msg := NewMessage(coap.NonConfirmable, coap.Content, utils.GenMessageID(), rv.Token, []byte(""))

				msg.SetOption(coap.LocationPath, rv.Path())
				format := rv.Option(coap.ContentFormat)
				if format != nil {
					msg.SetOption(coap.ContentFormat, format)
				}

				observe := rv.Option(coap.Observe).(uint32)
				if observe == 0 {
					if observable {
						if !HasObservation(route, request.FromAddr) {
							AddObservation(route, string(rv.Token), request.FromAddr, format)
							log.Printf("Register observing: %s\n", route)
						} else {
							UpdateObservation(route, string(rv.Token), request.FromAddr, format)
							log.Printf("Update observing: %s\n", route)
						}
						msg.AddOption(coap.Observe, 1)
					}
				} else if observe == 1 {
					RemoveObservation(route, request.FromAddr)
					log.Printf("Remove observing: %s\n", route)
				}
				// Send response (with or not Observe option)
				Send(l, request.FromAddr, *msg)
			} else {
				// NOT IMPLEMENTED
				rv.Code = coap.NotImplemented
				rv.Type = coap.Acknowledgement
				Send(l, request.FromAddr, rv)
			}
		default:
			// Route NOT FOUND
			rv.Code = coap.NotFound
			rv.Type = coap.Acknowledgement
			Send(l, request.FromAddr, rv)
		}
	}
}

func checkRoute(path string) (string, bool) {
	if res, ok := resources.GetByName(path); ok {
		return res.Name, res.Observable
	}
	return "", false
}
