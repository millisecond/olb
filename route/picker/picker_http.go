package picker

import (
	"net/http"
	"github.com/millisecond/olb/model"
)

// Picker selects a target from a list of targets.
type HTTPPicker func(targets []*model.Target, w http.ResponseWriter, req *http.Request) *model.Target

// Pickers contains the available Picker functions.
// Update config/load.go#load after updating.
var HTTPPickers = map[string]HTTPPicker{
	"cookie": CookieHTTPPicker,
	"rnd":    RndHTTPPicker,
	//"rr":     RrHTTPPicker,
}

// rndPicker picks a random target from the list of targets.
func RndHTTPPicker(targets []*model.Target, w http.ResponseWriter, req *http.Request) *model.Target {
	return RndPicker(targets)
}

//// rrPicker picks the next target from a list of targets using round-robin.
//func RrHTTPPicker(targets []*model.Target, w http.ResponseWriter, req *http.Request) *model.Target {
//	return rrPicker(targets)
//}
