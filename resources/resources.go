package resources

import (
	"github.com/qualiapps/subject/utils"
)

type (
	Res map[string]Resource

	Resource struct {
		Name       string `json:"name"`
		Observable bool   `json:"observable"`
		Payload    string `json:"payload"`
		Editable   bool   `json:"-"` // hidden field for json only
	}
)

var (
	Resources Res
	WellKnown = "/.well-known/core"
)

func init() {
	Resources = make(Res)
	AddResource(WellKnown, false, false)
}

/**
 * Retrieves all resources
 * @return - returns Resources reference
 */
func GetAll() *Res {
	return &Resources
}

func GetByName(name string) (Resource, bool) {
	if res, ok := Resources[name]; ok {
		return res, true
	}
	return Resource{}, false
}

// Adds a new client
// @param io.Reader params - json
// @return ([]byte, bool)
func AddResource(name string, observable, editable bool) (*Resource, bool) {
	if !utils.IsEmpty(name) && !HasResource(name) {
		r := Resource{Name: name, Observable: observable, Editable: editable}
		Resources[name] = r
		return &r, true
	}
	return nil, false
}

// Removes the client by Id
// @param string id - md5 hash of host:port
// @return string, bool
func DeleteResource(name string) (string, bool) {
	if HasResource(name) && IsEditable(name) {
		delete(Resources, name)
		return name, true
	}
	return "", false
}

func IsEditable(name string) bool {
	if res, ok := GetByName(name); ok {
		return res.Editable
	}
	return false
}

/**
* Checks an item
 * @param string id - client id
  * @return Confug, bool
*/
func HasResource(name string) bool {
	if _, ok := Resources[name]; ok {
		return true
	}
	return false
}
