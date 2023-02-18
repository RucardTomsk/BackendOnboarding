package service

import (
	"context"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/google/uuid"
)

type DivisionService struct {
	storage *postgresStorage.DivisionStorage
}

func NewDivisionService(
	storage *postgresStorage.DivisionStorage,
) *DivisionService {
	return &DivisionService{
		storage: storage,
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
		divisionsMas = append(divisionsMas, model.DivisionObject{Name: division.Name, Description: division.Description})
	}

	return divisionsMas, nil
}
