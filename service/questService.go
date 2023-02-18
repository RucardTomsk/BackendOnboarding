package service

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/storage/dao/neo4jRoles"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/google/uuid"
)

type QuestService struct {
	userStorage     *postgresStorage.UserStorage
	aboutStorage    *postgresStorage.AboutStorage
	roleStorage     *neo4jRoles.RolesStorage
	divisionStorage *postgresStorage.DivisionStorage
	questStorage    *postgresStorage.QuestStorage
}

func NewQuestService(
	userStorage *postgresStorage.UserStorage,
	aboutStorage *postgresStorage.AboutStorage,
	roleStorage *neo4jRoles.RolesStorage,
	divisionStorage *postgresStorage.DivisionStorage,
	questStorage *postgresStorage.QuestStorage) *QuestService {
	return &QuestService{
		userStorage:     userStorage,
		aboutStorage:    aboutStorage,
		roleStorage:     roleStorage,
		divisionStorage: divisionStorage,
		questStorage:    questStorage,
	}
}

func (s *QuestService) GetAllUserQuest(userID uuid.UUID) ([]model.UserQuestObject, *base.ServiceError) {
	divisionsGuid, err := s.roleStorage.GetDivision(userID)
	if err != nil {
		return nil, base.NewNeo4jReadError(err)
	}

	var userQuestMas []model.UserQuestObject

	for _, divGuid := range divisionsGuid {
		divID, err := uuid.Parse(divGuid)
		if err != nil {
			return nil, base.NewParseEnumError(err)
		}
		division, err := s.divisionStorage.Retrieve(divID, context.TODO())
		var questMas []model.OverQuestObject
		for _, quest := range division.Quests {
			questPre, err := s.questStorage.RetrieveQuest(quest.ID, context.TODO())
			if err != nil {
				return nil, base.NewPostgresReadError(err)
			}
			var stageMas []model.StageObject
			for _, stage := range questPre.Stages {
				stageMas = append(stageMas, model.StageObject{
					ID:          stage.ID,
					Name:        stage.Description,
					Description: stage.Description,
				})
			}
			questMas = append(questMas, model.OverQuestObject{
				ID:          quest.ID,
				Name:        quest.Name,
				Description: quest.Description,
				Stages:      stageMas,
			})
		}
		userQuestMas = append(userQuestMas, model.UserQuestObject{
			Division: model.DivisionObject{
				ID:          division.ID,
				Name:        division.Name,
				Description: division.Description,
			},
			Quests: questMas,
		})
	}
	return userQuestMas, nil
}
