package twixter

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Neo4j struct {
	driver neo4j.Driver
	log    *logrus.Logger
}

func NewNeo4j(config *viper.Viper, log *logrus.Logger) *Neo4j {
	driver, err := neo4j.NewDriver(
		config.GetString("neo4j.uri"),
		neo4j.BasicAuth(config.GetString("neo4j.user"), config.GetString("neo4j.pass"), ""),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to neo4j")
	}

	return &Neo4j{driver: driver, log: log}
}
