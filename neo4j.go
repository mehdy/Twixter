package twixter

import (
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4j struct {
	driver neo4j.Driver
}

func NewNeo4j() *Neo4j {
	driver, err := neo4j.NewDriver(
		os.Getenv("NEO4J_URI"),
		neo4j.BasicAuth(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PASS"), ""),
	)
	if err != nil {
		fmt.Printf("Failed to connect to neo4j: %s", err)
		return nil
	}

	return &Neo4j{driver: driver}
}
