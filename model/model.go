package model

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/millisecond/olb/config"
)

func Dynamo(config *config.Config) (*dynamo.DB, error) {
	sess, err := session.NewSession(config.AWSConfig.Generate())
	if err != nil {
		return nil, err
	}
	db := dynamo.New(sess, config.AWSConfig.Generate())
	return db, nil
}



