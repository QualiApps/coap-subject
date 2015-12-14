package main

import (
	"bytes"
	_ "encoding/json"
	"github.com/dustin/go-coap"
	"github.com/qualiapps/subject/resources"
	"log"
	"net"
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

func SendAck(l *net.UDPConn, from *net.UDPAddr, mid uint16) error {
	m := coap.Message{
		Type:      coap.Acknowledgement,
		Code:      0,
		MessageID: mid,
	}
	return coap.Transmit(l, from, m)
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

	msg := coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.Content,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   []byte(buf.String()),
	}

	msg.SetOption(coap.ContentFormat, coap.AppLinkFormat)
	msg.SetOption(coap.LocationPath, m.Path())

	err := coap.Transmit(l, from, msg)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return
	}
}
