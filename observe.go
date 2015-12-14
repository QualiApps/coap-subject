package main

import (
	_ "github.com/dustin/go-coap"
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
	observableList map[string][]*Observation
)

func NewObservation(addr *net.UDPAddr, token string, resource string) *Observation {
	return &Observation{
		Addr:        addr,
		Token:       token,
		Resource:    resource,
		NotifyCount: 0,
	}
}

type Observation struct {
	Addr        *net.UDPAddr
	Token       string
	Resource    string
	NotifyCount int
}

func AddObservation(resource, token string, addr *net.UDPAddr) {
	observableList[resource] = append(observableList[resource], NewObservation(addr, token, resource))
}

func HasObservation(resource string, addr *net.UDPAddr) bool {
	obs := observableList[resource]
	if obs == nil {
		return false
	}

	for _, o := range obs {
		if o.Addr.String() == addr.String() {
			return true
		}
	}
	return false
}

func RemoveObservation(resource string, addr *net.UDPAddr) {
	obs := observableList[resource]
	if obs == nil {
		return
	}

	for i, o := range obs {
		if o.Addr.String() == addr.String() {
			observableList[resource] = append(obs[:i], obs[i+1:]...)
			return
		}
	}
}
