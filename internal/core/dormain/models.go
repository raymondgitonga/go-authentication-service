package dormain

type AuthRequest struct {
	Email string `json:"email,validate:required,email"`
}
