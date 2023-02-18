package neo4jRoles

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/enum"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/RucardTomsk/BackendOnboarding/storage/driver"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	salt       = "nsfgnstg45s5fbnsfdg"
	signingKey = "qwerqwerGS#jjsS"
)

func encryptString(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type RolesStorage struct {
	driver          *driver.Neo4jDriver
	userStorage     *postgresStorage.UserStorage
	divisionStorage *postgresStorage.DivisionStorage
}

func NewRolesStorage(
	driver *driver.Neo4jDriver,
	userStorage *postgresStorage.UserStorage,
	divisionStorage *postgresStorage.DivisionStorage) *RolesStorage {
	return &RolesStorage{
		driver:          driver,
		userStorage:     userStorage,
		divisionStorage: divisionStorage,
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

	divisions, err := s.divisionStorage.Get()
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

	for _, division := range divisions {
		result, err := session.Run("MATCH (d:Division) WHERE d.guid = $guid RETURN d", map[string]interface{}{
			"guid": division.ID.String(),
		})

		if err != nil {
			return err
		}

		if !result.Next() {
			_, err := session.Run("CREATE (d:Division {guid: $guid})", map[string]interface{}{
				"guid": division.ID.String(),
			})

			if err != nil {
				return err
			}
		}
	}

	result, err := session.Run("MATCH (a:Admin) RETURN a", map[string]interface{}{})
	if err != nil {
		return err
	}

	if !result.Next() {
		_, err := session.Run("CREATE (a:Admin)", map[string]interface{}{})
		if err != nil {
			return err
		}
	}

	userAdmin, err := s.userStorage.RetrieveTo(model.AdminEmail, encryptString(model.AdminPassword), context.TODO())
	if err != nil {
		return err
	}
	result, err = session.Run("MATCH (a:Admin)-[]->(u:User) WHERE u.guid = $userGuid RETURN a,u", map[string]interface{}{
		"userGuid": userAdmin.ID.String(),
	})

	if err != nil {
		return err
	}

	if !result.Next() {
		_, err := session.Run("MATCH (a:Admin) MATCH (u:User) WHERE u.guid = $userGuid CREATE (a)-[:role {value: $valueRole}]->(u)", map[string]interface{}{
			"userGuid":  userAdmin.ID.String(),
			"valueRole": enum.ADMIN.String(),
		})

		if err != nil {
			return err
		}
	}

	if !result.Next() {

	}

	return nil
}

func (s RolesStorage) IssueRole(userID uuid.UUID, divisionID uuid.UUID, role enum.Roles) error {
	session := s.driver.GetSession()
	defer session.Close()

	_, err := session.Run("MATCH (u:User) WHERE u.guid = $guidUser MATCH (d:Division) WHERE d.guid = $guidDivision CREATE (u)-[:Role {value: $roleValue}]->(d)", map[string]interface{}{
		"guidUser":     userID.String(),
		"guidDivision": divisionID.String(),
		"roleValue":    role.String(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s RolesStorage) CheckAdminRole(userID uuid.UUID) (bool, error) {
	session := s.driver.GetSession()
	defer session.Close()

	result, err := session.Run("MATCH (a:Admin)-[]->(u:User) WHERE u.guid = $userGuid RETURN u", map[string]interface{}{
		"userGuid": userID.String(),
	})

	if err != nil {
		return false, err
	}

	if result.Next() {
		return true, nil
	}

	return false, nil
}

func (s RolesStorage) GetRole(userID uuid.UUID, divisionID uuid.UUID) (*enum.Roles, error) {
	session := s.driver.GetSession()
	defer session.Close()

	result, err := session.Run("MATCH (u:User)-[r]->(d:Division) WHERE u.guid = $userGuid, d.guid = $divisionGuid RETURN r.value", map[string]interface{}{
		"userGuid":     userID.String(),
		"divisionGuid": divisionID.String(),
	})

	if err != nil {
		return nil, err
	}

	var role enum.Roles
	if result.Next() {
		role = enum.ParseRoles(result.Record().Values[0].(string))
	}

	return &role, nil
}

func (s RolesStorage) GetDivision(userID uuid.UUID) ([]string, error) {
	session := s.driver.GetSession()
	defer session.Close()

	result, err := session.Run("MATCH (u:User)-[]->(d:Division) WHERE u.guid = $guid RETURN d.guid", map[string]interface{}{
		"guid": userID.String(),
	})

	if err != nil {
		return nil, err
	}

	var guidDivisionMas []string

	for result.Next() {
		guidDivisionMas = append(guidDivisionMas, result.Record().Values[0].(string))
	}

	return guidDivisionMas, nil
}
