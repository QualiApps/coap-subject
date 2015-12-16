package models

import (
	"bytes"
	"encoding/json"
	"github.com/qualiapps/subject/resources"
	"github.com/qualiapps/subject/utils"
	"io"
	"log"
)

// Retrieves all resources
// @return - returns bytes array
func GetAllResources() []byte {
	res, err := json.Marshal(resources.GetAll())
	if err != nil {
		log.Printf("ENCODE RESOURCES: Json encode error: %#v\n", err)
		return nil
	}
	return res
}

// Adds a new resource
// @param string name - res name
// @param bool observable
// @return (*Resource, []byte, bool)
func AddResource(name string, observable bool) (*resources.Resource, []byte, bool) {
	if res, ok := resources.AddResource(name, observable, true); ok {
		resp, err := json.Marshal(res)
		if err != nil {
			log.Printf("ENCODE RESOURCE: Json encode error: %#v\n", err)
			return nil, nil, false
		}
		return res, resp, true
	}
	return nil, nil, false
}

func InitEvent(name string, payload io.Reader) (*resources.Resource, bool) {
	if !utils.IsEmpty(name) {
		if res, ok := resources.GetByName(name); ok {
			buf := new(bytes.Buffer)
			buf.ReadFrom(payload)
			res.Payload = buf.String()
			return &res, true
		}
	}
	return nil, false
}

// Removes resource by name
// @param string name - res name
// @return bool
func DeleteResource(name string) (string, bool) {
	return resources.DeleteResource(name)
}
