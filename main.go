package nilpointer

// UserData represents user information
type UserData struct {
	ID    string
	Name  string
	Email string
}

// Database interface for data storage operations
type Database interface {
	GetUserByID(id string) (*UserData, error)
}

// Logger interface for logging operations
type Logger interface {
	Log(message string)
}

// Notifier sends notifications and has methods to be mocked
type Notifier interface {
	SendEmail(email string, message string) error
}

// UserService is the high-level service we'll be using
type UserService interface {
	GetUserDetails(userID string) (*UserData, error)
	GetUserLogger() Logger
	GetNotifier() Notifier
}

// UserProcessor processes user data
type UserProcessor struct {
	service UserService
}

// NewUserProcessor creates a new UserProcessor
func NewUserProcessor(service UserService) *UserProcessor {
	return &UserProcessor{
		service: service,
	}
}

// ProcessUserData gets user data and processes it
// This can cause a nil pointer dereference if GetNotifier() returns nil
func (p *UserProcessor) ProcessUserData(userID string) string {
	// Get user details
	userData, err := p.service.GetUserDetails(userID)
	if err != nil {
		return "Error getting user data"
	}

	// Get the logger
	logger := p.service.GetUserLogger()
	if logger != nil {
		logger.Log("Processing user: " + userData.ID)
	}

	// Get the notifier and use it
	notifier := p.service.GetNotifier()

	// This will cause a nil pointer dereference if notifier is nil
	err = notifier.SendEmail(userData.Email, "Your data has been processed")
	if err != nil {
		return "Error sending notification"
	}

	return "User data processed successfully"
}
