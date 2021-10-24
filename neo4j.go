package twixter

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/spf13/viper"
)

type Neo4j struct {
	driver neo4j.Driver
}

func NewNeo4j(config *viper.Viper) *Neo4j {
	driver, err := neo4j.NewDriver(
		config.GetString("neo4j.uri"),
		neo4j.BasicAuth(config.GetString("neo4j.user"), config.GetString("neo4j.pass"), ""),
	)
	if err != nil {
		fmt.Printf("Failed to connect to neo4j: %s", err)
		return nil
	}

	return &Neo4j{driver: driver}
}
