package main

import (
	"bytes"
	"github.com/dustin/go-coap"
	"github.com/qualiapps/subject/resources"
	"github.com/qualiapps/subject/utils"

	"log"
	"net"
	"strings"
	"time"
)

const (
	Net        = "udp"
	MaxBufSize = 1500
)

type (
	HostPort struct {
		Net     string
		Address string
	}
	CoapServer struct {
		ConnStr      HostPort
		Conn         *net.UDPConn
		observations obsList
	}
)

func NewCoapServer(port *string) *CoapServer {
	connStr := HostPort{Net, ":" + *port}
	return &CoapServer{
		ConnStr:      connStr,
		observations: make(obsList),
	}
}

func (c *CoapServer) Start() {
	sAddr, err := net.ResolveUDPAddr(c.ConnStr.Net, c.ConnStr.Address)
	if err != nil {
		log.Fatalln(err)
		return
	}

	c.Conn, err = net.ListenUDP(c.ConnStr.Net, sAddr)
	if err != nil {
		log.Fatalln(err)
	}

	h := Handle{}
	buf := make([]byte, MaxBufSize)
	for {
		nr, fromAddr, err := c.Conn.ReadFromUDP(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			log.Printf("Can't read from UDP: %#v\n", err)
		}
		tmp := make([]byte, nr)
		copy(tmp, buf)

		h.Data = tmp
		h.FromAddr = fromAddr

		go c.ProcessData(h)
	}

}

func (c *CoapServer) Event(name, payload string) {
	or := c.observations[name]
	if or != nil {
		for _, r := range or {
			c.SendNotification(r, coap.Content, payload)
		}
	}
}

func (c *CoapServer) DeregisterResource(name string) {
	c.SendDeregister(name, c.observations[name])
}

func (c *CoapServer) DeregisterResources() {
	for route, or := range c.observations {
		c.SendDeregister(route, or)
	}
}

func (c *CoapServer) SendDeregister(route string, or []*Observation) {
	if or != nil {
		for _, r := range or {
			c.SendNotification(r, coap.NotFound, "")
		}
		log.Printf("OK.........Observation %s was deleted.\n", route)
		delete(c.observations, route)
	}

}

func (c *CoapServer) SendNotification(r *Observation, code coap.COAPCode, payload string) {
	msg := NewMessage(coap.NonConfirmable, code, utils.GenMessageID(), []byte(r.Token), []byte(payload))

	msg.SetOption(coap.LocationPath, strings.Split(r.Resource, "/")[1:])

	if r.ContentFormat != nil {
		msg.SetOption(coap.ContentFormat, r.ContentFormat.(coap.MediaType))
	}
	r.NotifyCount++
	msg.AddOption(coap.Observe, r.NotifyCount)

	c.Send(r.Addr, *msg)
}

func (c *CoapServer) Send(addr *net.UDPAddr, request coap.Message) bool {
	err := coap.Transmit(c.Conn, addr, request)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return false
	}
	return true
}

func (c *CoapServer) SendAck(from *net.UDPAddr, mid uint16) bool {
	m := coap.Message{
		Type:      coap.Acknowledgement,
		Code:      0,
		MessageID: mid,
	}
	return c.Send(from, m)
}

func (c *CoapServer) Discovery(from *net.UDPAddr, m *coap.Message) {
	var buf bytes.Buffer
	for _, r := range *resources.GetAll() {
		buf.WriteString("<")
		buf.WriteString(r.Name)
		buf.WriteString(">")
		// @TODO add attribs
		buf.WriteString(",")
	}

	msg := NewMessage(coap.Acknowledgement, coap.Content, m.MessageID, m.Token, []byte(buf.String()))
	msg.SetOption(coap.ContentFormat, coap.AppLinkFormat)
	msg.SetOption(coap.LocationPath, m.Path())

	go c.Send(from, *msg)
}
