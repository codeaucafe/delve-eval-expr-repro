package nilpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService mocks the UserService interface
type MockUserService struct {
	mock.Mock
}

// GetUserDetails mocks the GetUserDetails method
func (m *MockUserService) GetUserDetails(userID string) (*UserData, error) {
	args := m.Called(userID)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*UserData), args.Error(1)
}

// GetUserLogger mocks the GetUserLogger method
func (m *MockUserService) GetUserLogger() Logger {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(Logger)
}

// GetNotifier mocks the GetNotifier method
func (m *MockUserService) GetNotifier() Notifier {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(Notifier)
}

// MockLogger mocks the Logger interface
type MockLogger struct {
	mock.Mock
}

// Log mocks the Log method
func (m *MockLogger) Log(message string) {
	m.Called(message)
}

// MockNotifier mocks the Notifier interface
type MockNotifier struct {
	mock.Mock
}

// SendEmail mocks the SendEmail method
func (m *MockNotifier) SendEmail(email string, message string) error {
	args := m.Called(email, message)
	return args.Error(0)
}

func TestUserProcessor_ProcessUserData_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockUserService)
	mockLogger := new(MockLogger)
	mockNotifier := new(MockNotifier)

	// Set up mock return values
	userData := &UserData{ID: "123", Name: "Test User", Email: "test@example.com"}
	mockService.On("GetUserDetails", "123").Return(userData, nil)
	mockService.On("GetUserLogger").Return(mockLogger)
	mockService.On("GetNotifier").Return(mockNotifier)

	// Set up logger expectation
	mockLogger.On("Log", "Processing user: 123").Return()

	// Set up notifier expectation
	mockNotifier.On("SendEmail", "test@example.com", "Your data has been processed").Return(nil)

	// Create the processor with the mock
	processor := NewUserProcessor(mockService)

	// Call the method
	result := processor.ProcessUserData("123")

	// Assert the result
	assert.Equal(t, "User data processed successfully", result)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}

func TestUserProcessor_ProcessUserData_NilPointerDereference(t *testing.T) {
	// Create mock service
	mockService := new(MockUserService)
	mockLogger := new(MockLogger)

	// Set up mock return values
	userData := &UserData{ID: "123", Name: "Test User", Email: "test@example.com"}
	mockService.On("GetUserDetails", "123").Return(userData, nil)
	mockService.On("GetUserLogger").Return(mockLogger)

	// Return nil for the notifier - this will cause a nil pointer dereference
	mockService.On("GetNotifier").Return(nil)

	// Set up logger expectation
	mockLogger.On("Log", "Processing user: 123").Return()

	// Create the processor with the mock
	processor := NewUserProcessor(mockService)

	// This will panic with a nil pointer dereference when it tries to call notifier.SendEmail()
	// Normally, we'd use a defer/recover to catch this, but for demonstration purposes we're
	// letting it panic to show the nil pointer dereference in delve
	result := processor.ProcessUserData("123")

	// This won't be reached
	assert.Equal(t, "User data processed successfully", result)
}
