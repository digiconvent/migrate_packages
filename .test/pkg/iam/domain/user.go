package iam_domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        *uuid.UUID // since this can be optional when creating a user
	Name      string
	Email     string
	Password  string
	Enabled   bool
	CreatedAt time.Time
}
