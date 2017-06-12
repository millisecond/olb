package route

import (
	"net/url"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/millisecond/olb/config"
	"github.com/millisecond/olb/metrics"
)

const HASHKEY_TARGET = "Target"

type Target struct {
	Type string

	ID string

	// Service is the name of the service the targetURL points to
	Service string

	// Tags are the list of tags for this target
	Tags []string

	// StripPath will be removed from the front of the outgoing
	// request path
	StripPath string

	// TLSSkipVerify disables certificate validation for upstream
	// TLS connections.
	TLSSkipVerify bool

	// Host signifies what the proxy will set the Host header to.
	// The proxy does not modify the Host header by default.
	// When Host is set to 'dst' the proxy will use the host name
	// of the target host for the outgoing request.
	Host string

	// URL is the endpoint the service instance listens on
	URL *url.URL

	// FixedWeight is the weight assigned to this target.
	// If the value is 0 the targets weight is dynamic.
	FixedWeight float64

	// Weight is the actual weight for this service in percent.
	Weight float64

	// Timer measures throughput and latency of this target
	Timer metrics.Timer

	// timerName is the name of the timer in the metrics registry
	timerName string
}

func (t *Target) put(config *config.Config) error {
	sess, err := session.NewSession(config.AWSConfig.Generate())
	if err != nil {
		return err
	}
	db := dynamo.New(sess, config.AWSConfig.Generate())
	table := db.Table(config.AWSConfig.DynamoTableName)
	//encoded := dynamo.AWSEncoding(t)
	put := table.Put(t)
	return put.Run()
}
