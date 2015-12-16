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

type Handle struct {
	Data     []byte
	FromAddr *net.UDPAddr
}

/**
 * Processing request
 * @param *net.UDPConn l - connection instance
 * @param Request request - res data
 */
func (c *CoapServer) ProcessData(h Handle) {
	// parse to CoAP struct
	rv := coap.Message{}
	err := rv.UnmarshalBinary(h.Data)

	if err == nil {
		path := string(os.PathSeparator) + strings.Join(rv.Path(), string(os.PathSeparator))
		route, observable := resources.CheckRoute(path)
		switch path {
		case resources.WellKnown:
			c.Discovery(h.FromAddr, &rv)
		case route:
			// if observe option
			if rv.IsObservable() {
				msg := NewMessage(coap.NonConfirmable, coap.Content, utils.GenMessageID(), rv.Token, []byte(""))

				msg.SetOption(coap.LocationPath, rv.Path())
				format := rv.Option(coap.ContentFormat)
				if format != nil {
					msg.SetOption(coap.ContentFormat, format)
				}

				// retrieves observe value
				observe := rv.Option(coap.Observe).(uint32)
				// observe = 0, register request
				if observe == 0 {
					// if resource is observable
					if observable {
						if !c.HasObservation(route, h.FromAddr) {
							c.AddObservation(route, string(rv.Token), h.FromAddr, format)
							log.Printf("Register observing: name-%s addr-%s\n", route, h.FromAddr.String())
						} else {
							c.UpdateObservation(route, string(rv.Token), h.FromAddr, format)
							log.Printf("Update observing: %s\n", route)
						}
						msg.AddOption(coap.Observe, 1)
					}
				} else if observe == 1 {
					c.RemoveObservation(route, h.FromAddr)
					log.Printf("Remove observing: name-%s addr-%s\n", route, h.FromAddr.String())
				}
				// Send response (with or not Observe option)
				c.Send(h.FromAddr, *msg)
			} else {
				// NOT IMPLEMENTED
				rv.Code = coap.NotImplemented
				rv.Type = coap.Acknowledgement
				c.Send(h.FromAddr, rv)
			}
		default:
			// Route NOT FOUND
			rv.Code = coap.NotFound
			rv.Type = coap.Acknowledgement
			c.Send(h.FromAddr, rv)
		}
	}
}
