package local

import (
	"context"
	"errors"
	"github.com/yerlanov/go-tour/common/session"
	"github.com/yerlanov/go-tour/main-api/internal/storage"
	"github.com/yerlanov/go-tour/main-api/internal/storage/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type service struct {
	userRepository user.Repository
	session        session.Session
}

type Service interface {
	Login(ctx context.Context, email, password string) (*user.User, error)
	Registration(ctx context.Context, dto UserCreateDto) (*string, error)
}

func NewService(storage storage.Storage, session session.Session) Service {
	return &service{
		userRepository: user.NewRepository(storage),
		session:        session,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, u.Password) {
		return nil, errors.New("invalid password")
	}

	sess, err := session.GetSessionFromContext(ctx)
	if err != nil {
		return nil, err
	}

	sess.Values["userId"] = u.ID.Hex()
	err = s.session.Set(ctx, sess)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Registration(ctx context.Context, dto UserCreateDto) (*string, error) {
	hash, err := HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	createUser := user.User{
		ID:        primitive.ObjectID{},
		Email:     dto.Email,
		Password:  hash,
		FirstName: dto.GivenName,
		LastName:  dto.FamilyName,
		Provider:  "local",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return s.userRepository.Upsert(ctx, createUser)
}
