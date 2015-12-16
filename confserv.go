package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/qualiapps/subject/controllers"
	"github.com/qualiapps/subject/resources"
	"log"
	"net/http"
	"strings"
)

func ServConfig(reg chan resources.Resource, rm chan string, ev chan resources.Resource) {
	// init router
	router := httprouter.New()

	register := func(r resources.Resource) {
		reg <- r
	}

	delete := func(name string) {
		rm <- name
	}

	event := func(r resources.Resource) {
		ev <- r
	}

	// Init controller
	controller := controllers.NewConfigController(register, delete, event)

	// Get resources list
	router.GET("/resources", controller.GetResources)

	// Add a new resource
	router.POST("/resources/add/:obs/*name", controller.AddResource)

	// Send Event
	router.POST("/resources/event/*name", controller.InitEvent)

	// Remove resource
	router.DELETE("/resources/delete/*name", controller.RemoveResource)

	log.Fatal(
		http.ListenAndServe(
			strings.Join([]string{"", *ConfigPort}, ":"),
			router,
		),
	)

}
