package main

import (
	"net"
)

type (
	obsList map[string][]*Observation
)

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

func (c *CoapServer) AddObservation(resource, token string, addr *net.UDPAddr, format interface{}) {
	c.observations[resource] = append(c.observations[resource], NewObservation(addr, format, token, resource))
}

func (c *CoapServer) UpdateObservation(resource, token string, addr *net.UDPAddr, format interface{}) {
	if obs, ok := c.observations[resource]; ok {
		for i, o := range obs {
			if o.Addr.String() == addr.String() {
				c.observations[resource][i].Token = token
				c.observations[resource][i].ContentFormat = format
				break
			}
		}
	}
}

func (c *CoapServer) HasObservation(resource string, addr *net.UDPAddr) bool {
	hasItem := false
	if obs, ok := c.observations[resource]; ok {
		for _, o := range obs {
			if o.Addr.String() == addr.String() {
				hasItem = true
				break
			}
		}
	}
	return hasItem
}

func (c *CoapServer) RemoveObservation(resource string, addr *net.UDPAddr) {
	if obs, ok := c.observations[resource]; ok {
		for i, o := range obs {
			if o.Addr.String() == addr.String() {
				c.observations[resource] = append(obs[:i], obs[i+1:]...)
				break
			}
		}
	}
}
