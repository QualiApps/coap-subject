package models

import (
	"bytes"
	"encoding/json"
	"github.com/qualiapps/subject/resources"
	"github.com/qualiapps/subject/utils"
	"io"
	"log"
)

/**
 * Retrieves all clients
 * @return - returns bytes array
 */
func GetAllResources() []byte {
	res, err := json.Marshal(resources.GetAll())
	if err != nil {
		log.Printf("ENCODE RESOURCES: Json encode error: %#v\n", err)
		return nil
	}
	return res
}

/**
 * Adds a new client
 * @param io.Reader params - json
 * @return ([]byte, bool)
 */
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

/**
 * Removes the client by Id
 * @param string id - md5 hash of host:port
 * @return bool
 */
func DeleteResource(name string) (string, bool) {
	return resources.DeleteResource(name)
}
