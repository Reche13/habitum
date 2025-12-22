package user

type CreateUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserPayload struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
