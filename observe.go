package main

import (
	"net"
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
	obsList map[string][]*Observation
)

var observableList = make(obsList)

func NewObservation(addr *net.UDPAddr, format interface{}, token, resource string) *Observation {
	return &Observation{
		Addr:          addr,
		Token:         token,
		Resource:      resource,
		ContentFormat: format,
		NotifyCount:   0,
	}
}

type Observation struct {
	Addr          *net.UDPAddr
	Token         string
	Resource      string
	ContentFormat interface{}
	NotifyCount   int
}

func AddObservation(resource, token string, addr *net.UDPAddr, format interface{}) {
	observableList[resource] = append(observableList[resource], NewObservation(addr, format, token, resource))
}

func UpdateObservation(resource, token string, addr *net.UDPAddr, format interface{}) {
	if obs, ok := observableList[resource]; ok {
		for i, o := range obs {
			if o.Addr.String() == addr.String() {
				observableList[resource][i].Token = token
				observableList[resource][i].ContentFormat = format
				break
			}
		}
	}
}

func HasObservation(resource string, addr *net.UDPAddr) bool {
	hasItem := false
	if obs, ok := observableList[resource]; ok {
		for _, o := range obs {
			if o.Addr.String() == addr.String() {
				hasItem = true
				break
			}
		}
	}
	return hasItem
}

func RemoveObservation(resource string, addr *net.UDPAddr) {
	if obs, ok := observableList[resource]; ok {
		for i, o := range obs {
			if o.Addr.String() == addr.String() {
				observableList[resource] = append(obs[:i], obs[i+1:]...)
				break
			}
		}
	}
}
