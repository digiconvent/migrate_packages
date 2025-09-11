package iam_user_repository

import (
	"github.com/digiconvent_clean_pkg/test/db"
	iam_domain "github.com/digiconvent_clean_pkg/test/pkg/iam/domain"
	"github.com/gofrs/uuid"
)

type IamUserRepositoryInterface interface {
	Create(user *iam_domain.User) (*uuid.UUID, error)
	Read(id *uuid.UUID) (*iam_domain.User, error)
	Update(user *iam_domain.User) error
	Delete(id *uuid.UUID) error
}

type IamUserRepository struct {
	db db.DatabaseInterface
}

func NewIamUserRepository(db db.DatabaseInterface) IamUserRepositoryInterface {
	return &IamUserRepository{
		db: db,
	}
}
