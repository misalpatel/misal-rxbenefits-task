package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rxbenefits/go-hw/internal/repository"
)

func TestNewFilmRepository(t *testing.T) {
	assert.NotNil(t, repository.NewFilmRepository)
}

func TestNewCommentRepository(t *testing.T) {
	assert.NotNil(t, repository.NewCommentRepository)
}
