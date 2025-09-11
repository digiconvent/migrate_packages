package iam_repository

import iam_user_repository "github.com/digiconvent_clean_pkg/test/pkg/iam/repository/user"

type IamRepository struct {
	UserRepository iam_user_repository.IamUserRepositoryInterface
	// repositories for other entities here
}

func NewIamRepository(db DatabaseInterface) *IamRepository {
	return &IamRepository{
		UserRepository: iam_user_repository.NewIamUserRepository(db),
		// repositories for other entities here
	}
}
