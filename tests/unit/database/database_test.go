package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rxbenefits/go-hw/internal/database"
)

func TestWithDBHost(t *testing.T) {
	// Test that WithDBHost function exists and can be called
	withHost := database.WithDBHost("test-host")
	assert.NotNil(t, withHost)
}

func TestWithDBPort(t *testing.T) {
	// Test that WithDBPort function exists and can be called
	withPort := database.WithDBPort("5433")
	assert.NotNil(t, withPort)
}

func TestWithDBUser(t *testing.T) {
	// Test that WithDBUser function exists and can be called
	withUser := database.WithDBUser("test-user")
	assert.NotNil(t, withUser)
}

func TestWithDBPassword(t *testing.T) {
	// Test that WithDBPassword function exists and can be called
	withPassword := database.WithDBPassword("test-password")
	assert.NotNil(t, withPassword)
}

func TestWithDBName(t *testing.T) {
	// Test that WithDBName function exists and can be called
	withDBName := database.WithDBName("dvdrental_test")
	assert.NotNil(t, withDBName)
}

func TestInitDB_WithOptions(t *testing.T) {
	// Test with custom options
	db, err := database.InitDB(
		database.WithDBHost("test-host"),
		database.WithDBPort("5433"),
		database.WithDBUser("test-user"),
		database.WithDBPassword("test-password"),
		database.WithDBName("dvdrental_test"),
	)

	// Should fail because we can't connect to test database
	require.Error(t, err)
	assert.Nil(t, db)
}

func TestInitDB_DefaultOptions(t *testing.T) {
	// Test with default options
	db, err := database.InitDB()

	// Should fail because we can't connect to database
	require.Error(t, err)
	assert.Nil(t, db)
}
