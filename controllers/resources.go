package controllers

import (
	_ "encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/qualiapps/subject/models"
	"github.com/qualiapps/subject/resources"
	"net/http"
	"strconv"
)

type (
	// ConfController represents the controller for operating on the Config resource
	ConfigController struct {
		OnRegister func(c resources.Resource)
		OnEvent    func(c resources.Resource)
		OnDelete   func(name string)
	}
)

func NewConfigController(reg func(r resources.Resource), rm func(name string), ev func(r resources.Resource)) *ConfigController {
	return &ConfigController{reg, ev, rm}
}

/**
 * Retrieves all clients
 * @return json
 */
func (c *ConfigController) GetResources(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res := models.GetAllResources()

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Fprintf(w, "%s", res)
}

func (c *ConfigController) InitEvent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	status := 404
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	if resource, ok := models.InitEvent(p.ByName("name"), r.Body); ok {
		c.Event(*resource)
		status = 200
	}
	w.WriteHeader(status)
}

/**
 * Adds a new client
 */
func (c *ConfigController) AddResource(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		obs bool
		err error
	)
	if obs, err = strconv.ParseBool(p.ByName("obs")); err != nil {
		obs = false
	}
	resource, response, ok := models.AddResource(p.ByName("name"), obs)
	if ok {
		c.Register(*resource)
	}

	// Write content-type, statuscode, response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", response)
}

/**
 * Removes an existing client
 */
func (c *ConfigController) RemoveResource(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	status := 400
	if name, ok := models.DeleteResource(p.ByName("name")); ok {
		status = 200
		c.Delete(name)
	}
	w.WriteHeader(status)
}

func (c *ConfigController) Register(res resources.Resource) {
	c.OnRegister(res)
}

func (c *ConfigController) Event(res resources.Resource) {
	c.OnEvent(res)
}

func (c *ConfigController) Delete(name string) {
	c.OnDelete(name)
}
