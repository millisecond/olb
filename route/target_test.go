package route

import (
	"github.com/millisecond/olb/config"
	"github.com/millisecond/olb/uuid"
	"testing"
)

func TestTargetPut(t *testing.T) {
	target := &Target{ID: uuid.NewUUID()}
	err := target.put(&config.Config{
		AWSConfig: config.AWSConfig{
			Region:          "us-west2",
			Endpoint:        "http://localhost:8000",
			DynamoTableName: "OLB",
			Key: "key",
			Secret: "secret",
		},
	})
	if err != nil {
		t.Errorf("Got an error: ", err)
	}
}
