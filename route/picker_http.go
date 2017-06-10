package route


import (
	"net/http"
)

const HTTP_COOKIE_CONTEXT_KEY="HTTP_COOKIE_CONTEXT_KEY"
const HTTP_COOKIENAME="OLBKEY"

// picker selects a target from a list of targets.
type httpPicker func(r *Route, req *http.Request) *Target

// Pickers contains the available picker functions.
// Update config/load.go#load after updating.
var HTTPPickers = map[string]httpPicker{
	"cookie": cookieHTTPPicker,
	"rnd": rndHTTPPicker,
	"rr":  rrHTTPPicker,
}

// rndPicker picks a random target from the list of targets.
func cookieHTTPPicker(r *Route, req *http.Request) *Target {
	return rndPicker(r)
}

// rndPicker picks a random target from the list of targets.
func rndHTTPPicker(r *Route, req *http.Request) *Target {
	return rndPicker(r)
}

// rrPicker picks the next target from a list of targets using round-robin.
func rrHTTPPicker(r *Route, req *http.Request) *Target {
	return rrPicker(r)
}
