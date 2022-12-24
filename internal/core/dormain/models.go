package dormain

type AuthRequest struct {
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
	Token  string `json:"token,omitempty"`
}
