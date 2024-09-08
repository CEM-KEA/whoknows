package unit_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockViper simulates the ViperConfig interface
type MockViper struct {
	mock.Mock
}

func (m *MockViper) SetConfigFile(file string) {
	m.Called(file)
}

func (m *MockViper) ReadInConfig() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockViper) SetConfigName(name string) {
	m.Called(name)
}

func (m *MockViper) SetConfigType(typ string) {
	m.Called(typ)
}

func (m *MockViper) AddConfigPath(path string) {
	m.Called(path)
}

func (m *MockViper) AutomaticEnv() {
	m.Called()
}

func (m *MockViper) SetEnvPrefix(prefix string) {
	m.Called(prefix)
}

func (m *MockViper) SetEnvKeyReplacer(r *strings.Replacer) {
	m.Called(r)
}

func (m *MockViper) SetDefault(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockViper) Unmarshal(rawVal interface{}) error {
	args := m.Called(rawVal)
	return args.Error(0)
}

// Test LoadEnv function
func TestLoadEnvSuccess(t *testing.T) {
	mockViper := new(MockViper)
	mockViper.On("SetConfigFile", ".env").Return()
	mockViper.On("ReadInConfig").Return(nil)

	// Call LoadEnv with mock
	err := config.LoadEnv(mockViper)

	assert.NoError(t, err)
	mockViper.AssertExpectations(t)
}

func TestLoadEnvFailure(t *testing.T) {
	mockViper := new(MockViper)
	mockViper.On("SetConfigFile", ".env").Return()
	mockViper.On("ReadInConfig").Return(errors.New("file not found"))

	err := config.LoadEnv(mockViper)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading .env file")
}

// Test LoadYaml function
func TestLoadYamlSuccess(t *testing.T) {
	mockViper := new(MockViper)
	mockViper.On("SetConfigName", "config").Return()
	mockViper.On("SetConfigType", "yaml").Return()
	mockViper.On("AddConfigPath", "./internal/config").Return()
	mockViper.On("ReadInConfig").Return(nil)

	err := config.LoadYaml(mockViper)

	assert.NoError(t, err)
	mockViper.AssertExpectations(t)
}

func TestLoadYamlFailure(t *testing.T) {
	mockViper := new(MockViper)
	mockViper.On("SetConfigName", "config").Return()
	mockViper.On("SetConfigType", "yaml").Return()
	mockViper.On("AddConfigPath", "./internal/config").Return()
	mockViper.On("ReadInConfig").Return(errors.New("config file not found"))

	err := config.LoadYaml(mockViper)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config.yaml file")
}

// Test LoadConfig function
func TestLoadConfigSuccess(t *testing.T) {
	mockViper := new(MockViper)

	// Mocking the .env file load
	mockViper.On("SetConfigFile", ".env").Return()
	mockViper.On("ReadInConfig").Return(nil)

	// Mocking the config.yaml file load
	mockViper.On("SetConfigName", "config").Return()
	mockViper.On("SetConfigType", "yaml").Return()
	mockViper.On("AddConfigPath", "./internal/config").Return()
	mockViper.On("ReadInConfig").Return(nil)

	// Mock environment variable reading and other config setups
	mockViper.On("AutomaticEnv").Return()
	mockViper.On("SetEnvPrefix", "API").Return()
	mockViper.On("SetEnvKeyReplacer", mock.Anything).Return()
	mockViper.On("SetDefault", "server.port", 8080).Return()
	mockViper.On("Unmarshal", mock.Anything).Return(nil)

	// Call LoadConfig with the mockViper
	err := config.LoadConfig(mockViper)

	// Assert no errors occurred
	assert.NoError(t, err)

	// Ensure all mock expectations were met
	mockViper.AssertExpectations(t)
}

func TestLoadConfigFailure(t *testing.T) {
	mockViper := new(MockViper)
	mockViper.On("SetConfigFile", ".env").Return()
	mockViper.On("ReadInConfig").Return(errors.New("file not found"))

	err := config.LoadConfig(mockViper)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading .env file")
}
