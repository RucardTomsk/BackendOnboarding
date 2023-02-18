package base

import "github.com/google/uuid"

type Blame string

const (
	BlameUser     Blame = "User"
	BlamePostgres Blame = "Postgres"
	BlameServer   Blame = "Server"
	BlameUnknown  Blame = "Unknown"
	BlameNeo4j    Blame = "Neo4j"
)

// ResponseOK is a base OK response from server.
type ResponseOK struct {
	Status     string `json:"status" example:"OK"`
	TrackingID string `json:"trackingID" example:"12345678-1234-1234-1234-000000000000"`
}

// ResponseOKWithGUID is a base OK response from server with additional GUID in answer.
type ResponseOKWithGUID struct {
	Status     string    `json:"status" example:"OK"`
	TrackingID string    `json:"trackingID" example:"12345678-1234-1234-1234-000000000000"`
	GUID       uuid.UUID `json:"GUID" example:"12345678-1234-1234-1234-000000000000"`
}

// ResponseFailure is a general error response from server.
type ResponseFailure struct {
	Status     string `json:"status" example:"Error"`
	Blame      Blame  `json:"blame" example:"Guilty System"`
	TrackingID string `json:"trackingID" example:"12345678-1234-1234-1234-000000000000"`
	Message    string `json:"message" example:"error occurred"`
}

type ResponseOKWithJWT struct {
	Status     string    `json:"status" example:"OK"`
	TrackingID string    `json:"trackingID" example:"12345678-1234-1234-1234-000000000000"`
	GUID       uuid.UUID `json:"GUID" example:"12345678-1234-1234-1234-000000000000"`
	JWT        string    `json:"JWT"`
}
