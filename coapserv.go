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

func ServCoap(listener chan *net.UDPConn, handler chan Request, conStr HostPort) {
	sAddr, err := net.ResolveUDPAddr(conStr.Net, conStr.Address)
	if err != nil {
		log.Fatalln(err)
		return
	}

	l, err := net.ListenUDP(conStr.Net, sAddr)
	if err != nil {
		log.Fatalln(err)
	}

	listener <- l

	buf := make([]byte, MaxBufSize)
	request := Request{}
	for {
		nr, fromAddr, err := l.ReadFromUDP(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			log.Printf("Can't read from UDP: %#v\n", err)
		}
		tmp := make([]byte, nr)
		copy(tmp, buf)

		request.Data = tmp
		request.FromAddr = fromAddr

		handler <- request
	}

}

func SendNotification(l *net.UDPConn, r *Observation, code coap.COAPCode, payload string) {
	msg := NewMessage(coap.NonConfirmable, code, utils.GenMessageID(), []byte(r.Token), []byte(payload))

	msg.SetOption(coap.LocationPath, strings.Split(r.Resource, "/")[1:])

	if r.ContentFormat != nil {
		msg.SetOption(coap.ContentFormat, r.ContentFormat.(coap.MediaType))
	}
	r.NotifyCount++
	msg.AddOption(coap.Observe, r.NotifyCount)

	go Send(l, r.Addr, *msg)
}

func Send(l *net.UDPConn, addr *net.UDPAddr, request coap.Message) bool {
	err := coap.Transmit(l, addr, request)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return false
	}
	return true
}

func SendAck(l *net.UDPConn, from *net.UDPAddr, mid uint16) bool {
	m := coap.Message{
		Type:      coap.Acknowledgement,
		Code:      0,
		MessageID: mid,
	}
	return Send(l, from, m)
}

func Discovery(l *net.UDPConn, from *net.UDPAddr, m *coap.Message) {
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

	go Send(l, from, *msg)
}
