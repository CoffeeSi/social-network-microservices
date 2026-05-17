package event

type UserVerificationEvent struct {
	EventType string `json:"event_type"`
	OccuredAt string `json:"occured_at"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}
