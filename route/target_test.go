package route

import (
	"github.com/millisecond/olb/config"
	"github.com/millisecond/olb/install"
	"github.com/millisecond/olb/uuid"
	"testing"
)

func TestTargetPut(t *testing.T) {
	target := &Target{
		Type: HASHKEY_TARGET,
		ID:   uuid.NewUUID(),
	}
	cfg := &config.Config{
		AWSConfig: config.AWSConfig{
			Region:          "us-west2",
			Endpoint:        "http://localhost:8000",
			DynamoTableName: "OLB",
			Key:             "key",
			Secret:          "secret",
		},
	}
	install.Install(cfg)
	err := target.put(cfg)
	if err != nil {
		t.Errorf("Got an error: ", err)
	}
}
