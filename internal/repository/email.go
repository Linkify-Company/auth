package repository

import (
	"auth/internal/repository/memory"
	"context"
	"github.com/Linkify-Company/common_utils/errify"
)

type EmailRepos struct {
	emailMemory *memory.EmailMemory
}

func NewEmailRepos() Email {
	return &EmailRepos{
		emailMemory: memory.NewEmailMemory(),
	}
}

func (m *EmailRepos) IsExist(ctx context.Context, email string) bool {
	return m.emailMemory.IsExist(email)
}

func (m *EmailRepos) Set(ctx context.Context, email string, code int) errify.IError {
	m.emailMemory.Set(email, code)
	return nil
}

func (m *EmailRepos) IsValid(ctx context.Context, email string, code int) bool {
	return m.emailMemory.IsValid(email, code)
}
