package event

const (
	UserRegisteredEventType   = "user.registered"
	UserVerificationEventType = "user.verification"
)

type UserRegisteredEvent struct {
	EventType  string `json:"event_type"`
	OccurredAt string `json:"occurred_at"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	DOB        string `json:"dob"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
type UserVerificationEvent struct {
	EventType string `json:"event_type"`
	OccuredAt string `json:"occured_at"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}
