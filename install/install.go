package install

import (
	"github.com/millisecond/olb/config"
	"github.com/millisecond/olb/model"
)

type DynamoTableDefinition struct {
	Type string `dynamo:"Type,hash"`
	ID   string `dynamo:"ID,range"`
}

func Install(cfg *config.Config) error {
	// generate cookie secret(s)
	// validate AWS creds

	// create dynamo table
	db, err := model.Dynamo(cfg)
	if err != nil {
		return err
	}

	return db.CreateTable(cfg.AWSConfig.DynamoTableName, &DynamoTableDefinition{}).Run()
}
