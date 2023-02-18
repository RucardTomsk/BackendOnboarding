package service

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/enum"
	"github.com/RucardTomsk/BackendOnboarding/storage/dao/neo4jRoles"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/google/uuid"
)

type DivisionService struct {
	storage      *postgresStorage.DivisionStorage
	roleStorage  *neo4jRoles.RolesStorage
	questStorage *postgresStorage.QuestStorage
}

func NewDivisionService(
	storage *postgresStorage.DivisionStorage,
	roleStorage *neo4jRoles.RolesStorage,
	questStorage *postgresStorage.QuestStorage,
) *DivisionService {
	return &DivisionService{
		storage:      storage,
		roleStorage:  roleStorage,
		questStorage: questStorage,
	}
}

func (s *DivisionService) Create(request *model.CreateDivisionRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {
	division := entity.Division{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := s.storage.Create(&division, context.TODO()); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	return &division.ID, nil
}

func (s *DivisionService) Get(ctx context.Context) ([]model.DivisionObject, *base.ServiceError) {
	divisions, err := s.storage.Get()
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	var divisionsMas []model.DivisionObject

	for _, division := range divisions {
		divisionsMas = append(divisionsMas, model.DivisionObject{
			ID:          division.ID,
			Name:        division.Name,
			Description: division.Description,
		})
	}

	return divisionsMas, nil
}

func (s *DivisionService) AddUser(userID uuid.UUID, divisionID uuid.UUID, role string, ctx context.Context) *base.ServiceError {
	if err := s.roleStorage.IssueRole(userID, divisionID, enum.ParseRoles(role)); err != nil {
		return base.NewNeo4jWriteError(err)
	}

	return nil
}

func (s *DivisionService) AddQuest(divisionID uuid.UUID, request *model.CreateQuestRequest, ctx context.Context) *base.ServiceError {
	quest := entity.Quest{
		Name:        request.Name,
		Description: request.Description,
		IsDone:      false,
		DivisionID:  divisionID,
	}

	if err := s.questStorage.CreateQuest(&quest, context.TODO()); err != nil {
		return base.NewPostgresWriteError(err)
	}

	division, err := s.storage.Retrieve(divisionID, context.TODO())
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	division.Quests = append(division.Quests, quest)

	if err := s.storage.Update(division, context.TODO()); err != nil {
		return base.NewPostgresWriteError(err)
	}

	return nil
}

func (s DivisionService) AddStage(questID uuid.UUID, request *model.CreateStageRequest, ctx context.Context) *base.ServiceError {
	stage := entity.Stage{
		Name:        request.Name,
		Description: request.Description,
		IsDone:      false,
		QuestID:     questID,
	}

	if err := s.questStorage.CreateStage(&stage, context.TODO()); err != nil {
		return base.NewPostgresWriteError(err)
	}

	quest, err := s.questStorage.RetrieveQuest(questID, context.TODO())
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	quest.Stages = append(quest.Stages, stage)

	if err := s.questStorage.UpdateQuest(quest, context.TODO()); err != nil {
		return base.NewPostgresWriteError(err)
	}

	return nil
}

func (s DivisionService) GetAllQuest() ([]model.QuestObject, *base.ServiceError) {
	quests, err := s.questStorage.GetQuest()
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	var questsMas []model.QuestObject

	for _, quest := range quests {
		var stagesMas []model.StageObject
		for _, stage := range quest.Stages {
			stagesMas = append(stagesMas, model.StageObject{
				ID:          stage.ID,
				Name:        stage.Name,
				Description: stage.Description,
			})
		}
		questsMas = append(questsMas, model.QuestObject{
			ID:          quest.ID,
			Name:        quest.Name,
			Description: quest.Description,
		})
	}

	return questsMas, nil
}
