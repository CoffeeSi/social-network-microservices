package event

const (
	UserRegisteredEventType = "user.registered"
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
