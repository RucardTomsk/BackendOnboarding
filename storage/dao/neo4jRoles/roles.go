package neo4jRoles

import (
	"errors"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/enum"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/RucardTomsk/BackendOnboarding/storage/driver"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolesStorage struct {
	driver      *driver.Neo4jDriver
	userStorage *postgresStorage.UserStorage
}

func NewRolesStorage(driver *driver.Neo4jDriver,
	userStorage *postgresStorage.UserStorage) *RolesStorage {
	return &RolesStorage{
		driver:      driver,
		userStorage: userStorage,
	}
}

func (s RolesStorage) CreateUser(userID uuid.UUID) error {
	session := s.driver.GetSession()
	defer session.Close()

	_, err := session.Run("CREATE (n:User {guid: $guid})", map[string]interface{}{
		"guid": userID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s RolesStorage) Migrations() error {
	users, err := s.userStorage.Get()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	session := s.driver.GetSession()
	defer session.Close()

	for _, user := range users {
		result, err := session.Run("MATCH (u:User) WHERE u.guid = $guid RETURN u", map[string]interface{}{
			"guid": user.ID.String(),
		})

		if err != nil {
			return err
		}

		if !result.Next() {
			_, err := session.Run("CREATE (u:User {guid: $guid})", map[string]interface{}{
				"guid": user.ID.String(),
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s RolesStorage) IssueRole(userID uuid.UUID, divisionID uuid.UUID, role enum.Roles) error {
	session := s.driver.GetSession()
	defer session.Close()

	_, err := session.Run("MATCH (u:User) WHERE u.guid = $guidUser MATCH (d:Division) WHERE d.guid = $guidDivision CREATE (u)-[:$roleValue]->(d)", map[string]interface{}{
		"guidUser":     userID,
		"guidDivision": divisionID,
		"roleValue":    role.String(),
	})

	if err != nil {
		return err
	}

	return nil
}
