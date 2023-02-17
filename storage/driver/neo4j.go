package driver

import (
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/internal/common"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jDriver struct {
	driver       neo4j.Driver
	DataBaseName string
}

func NewNeo4jDriver(config *common.Neo4jConfig) (*Neo4jDriver, error) {
	driver, err := neo4j.NewDriver(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		neo4j.BasicAuth(config.User, config.Password, ""),
	)

	if err != nil {
		return nil, err
	}
	if err := driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &Neo4jDriver{
		driver:       driver,
		DataBaseName: config.DataBaseName,
	}, nil
}

func (d *Neo4jDriver) GetSession() neo4j.Session {
	return d.driver.NewSession(neo4j.SessionConfig{
		DatabaseName: d.DataBaseName,
	})
}
