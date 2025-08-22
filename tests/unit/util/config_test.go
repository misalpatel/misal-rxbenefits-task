package util_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rxbenefits/go-hw/internal/util"
)

func TestInitConfig(t *testing.T) {
	// Test default values
	config := util.InitConfig()

	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "5432", config.DBPort)
	assert.Equal(t, "postgres", config.DBUser)
	assert.Equal(t, "postgres", config.DBPassword)
	assert.Equal(t, "dvdrental", config.DBName)
}

func TestInitConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	t.Setenv("DB_HOST", "test-host")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "test-user")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_NAME", "test-db")

	// Clean up after test
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	config := util.InitConfig()

	assert.Equal(t, "test-host", config.DBHost)
	assert.Equal(t, "5433", config.DBPort)
	assert.Equal(t, "test-user", config.DBUser)
	assert.Equal(t, "test-password", config.DBPassword)
	assert.Equal(t, "test-db", config.DBName)
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	t.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	value := util.GetEnv("TEST_VAR", "default")
	assert.Equal(t, "test-value", value)

	// Test with non-existing environment variable
	value = util.GetEnv("NON_EXISTENT_VAR", "default-value")
	assert.Equal(t, "default-value", value)
}
