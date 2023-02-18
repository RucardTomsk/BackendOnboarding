package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

const (
	salt       = "nsfgnstg45s5fbnsfdg"
	signingKey = "qwerqwerGS#jjsS"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_guid"`
}

type UserService struct {
	storage *postgresStorage.UserStorage
}

func NewUserService(
	storage *postgresStorage.UserStorage,
) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) Create(request *model.CreateUserRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {
	user := entity.User{
		Password: encryptString(request.Password),
		Email:    request.Email,
		About:    &entity.About{},
	}

	if err := s.storage.Create(&user, context.TODO()); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	return &user.ID, nil
}

func (s *UserService) UpdateAbout(request *model.UpdateAbout, userID uuid.UUID, ctx context.Context) error {
	user, err := s.storage.Retrieve(userID, context.TODO())
	if err != nil {
		return err
	}

	about := entity.About{
		Description: request.Description,
		FIO:         request.FIO,
		Contact:     request.Contact,
	}

	user.About = &about

	if err := s.storage.Update(user, context.TODO()); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GenerateToken(request *model.GenerateTokenRequest, ctx context.Context) (*model.Token, *base.ServiceError) {
	user, err := s.storage.RetrieveTo(request.Email, encryptString(request.Password), context.TODO())
	if err != nil {
		return nil, base.NewGenerateJWTError(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID.String(),
	})

	valueToken, err := token.SignedString([]byte(signingKey))
	return &model.Token{Value: valueToken}, nil
}

func (s *UserService) GetEmail(userID uuid.UUID, ctx context.Context) (string, *base.ServiceError) {
	user, err := s.storage.Retrieve(userID, context.TODO())
	if err != nil {
		return "", base.NewPostgresReadError(err)
	}

	return user.Email, nil
}

func (s *UserService) Get(ctx context.Context) ([]entity.User, *base.ServiceError) {
	users, err := s.storage.Get()
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	return users, nil
}

func (s *UserService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, base.NewParseJWTError(errors.New("invalid signing method"))
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", base.NewParseJWTError(errors.New("token claims are not of type *tokenClaims"))
	}

	return claims.UserID, nil
}

func encryptString(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
