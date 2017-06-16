package route

import (
	"log"
	"path"
	"strings"
	"github.com/millisecond/olb/model"
)

// matcher determines whether a host/path matches a route
type matcher func(uri string, r *model.Route) bool

// Matcher contains the available matcher functions.
// Update config/load.go#load after updating.
var Matcher = map[string]matcher{
	"prefix": prefixMatcher,
	"glob":   globMatcher,
}

// prefixMatcher matches path to the routes' path.
func prefixMatcher(uri string, r *model.Route) bool {
	return strings.HasPrefix(uri, r.Path)
}

// globMatcher matches path to the routes' path using globbing.
func globMatcher(uri string, r *model.Route) bool {
	var hasMatch, err = path.Match(r.Path, uri)
	if err != nil {
		log.Printf("[ERROR] Glob matching error %s for path %s route %s", err, uri, r.Path)
		return false
	}
	return hasMatch
}
