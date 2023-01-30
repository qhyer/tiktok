package neo4j

import (
	"tiktok/pkg/constants"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

func Init() {
	result, err := neo4j.NewDriverWithContext(constants.Neo4jDefaultURI, neo4j.BasicAuth(constants.Neo4jUser, constants.Neo4jPassword, ""))
	if err != nil {
		panic(err)
	}
	driver = result
}
