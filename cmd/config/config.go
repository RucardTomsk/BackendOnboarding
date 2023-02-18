package config

import "github.com/RucardTomsk/BackendOnboarding/internal/common"

type Config struct {
	Postgres common.DatabaseConfig
	Server   common.ServerConfig
	Neo4j    common.Neo4jConfig
}
