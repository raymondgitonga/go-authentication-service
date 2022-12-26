package dormain

type AuthResponse struct {
	Name   string `json:"name,omitempty"`
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
	Token  string `json:"token,omitempty"`
}
