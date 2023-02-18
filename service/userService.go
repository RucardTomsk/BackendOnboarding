package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/entity"
	"github.com/RucardTomsk/BackendOnboarding/storage/dao/neo4jRoles"
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
	userStorage     *postgresStorage.UserStorage
	aboutStorage    *postgresStorage.AboutStorage
	roleStorage     *neo4jRoles.RolesStorage
	divisionStorage *postgresStorage.DivisionStorage
}

func NewUserService(
	userStorage *postgresStorage.UserStorage,
	aboutStorage *postgresStorage.AboutStorage,
	roleStorage *neo4jRoles.RolesStorage,
	divisionStorage *postgresStorage.DivisionStorage,
) *UserService {
	return &UserService{
		userStorage:     userStorage,
		aboutStorage:    aboutStorage,
		roleStorage:     roleStorage,
		divisionStorage: divisionStorage,
	}
}

func (s *UserService) Create(request *model.CreateUserRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {

	user := entity.User{
		Password: encryptString(request.Password),
		Email:    request.Email,
	}

	if err := s.userStorage.Create(&user, context.TODO()); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	about := entity.About{
		UserID: user.ID,
	}

	if err := s.aboutStorage.Create(&about, context.TODO()); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	user.About = &about
	user.AboutID = about.ID

	if err := s.userStorage.Update(&user, context.TODO()); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	return &user.ID, nil
}

func (s *UserService) UpdateAbout(request *model.UpdateAbout, userID uuid.UUID, ctx context.Context) *base.ServiceError {
	user, err := s.userStorage.Retrieve(userID, context.TODO())
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	user.About.FIO = request.FIO
	user.About.Description = request.Description
	user.About.Contact = request.Contact

	if err := s.aboutStorage.Update(user.About, context.TODO()); err != nil {
		return base.NewPostgresWriteError(err)
	}
	//if err := s.userStorage.Update(user, context.TODO()); err != nil {
	//return base.NewPostgresWriteError(err)
	//}

	return nil
}

func (s *UserService) GenerateToken(request *model.GenerateTokenRequest, ctx context.Context) (*model.Token, *base.ServiceError) {
	user, err := s.userStorage.RetrieveTo(request.Email, encryptString(request.Password), context.TODO())
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

func (s *UserService) GetInfo(userID uuid.UUID, ctx context.Context) (*model.UserObject, *base.ServiceError) {
	user, err := s.userStorage.Retrieve(userID, context.TODO())
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	return &model.UserObject{
		ID:    user.ID,
		Email: user.Email,
		About: model.AboutObject{
			FIO:         user.About.FIO,
			Description: user.About.Description,
			Contact:     user.About.Contact,
		},
	}, nil
}

func (s *UserService) Get(ctx context.Context) ([]model.UserObject, *base.ServiceError) {
	users, err := s.userStorage.Get()
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	var userMas []model.UserObject

	for _, user := range users {
		userMas = append(userMas, model.UserObject{
			ID:    user.ID,
			Email: user.Email,
			About: model.AboutObject{
				FIO:         user.About.FIO,
				Description: user.About.Description,
				Contact:     user.About.Contact,
			},
		})
	}
	return userMas, nil
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
