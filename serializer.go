package gota

type ListResponse[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

type LoginResp struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresIn  int       `json:"access_expires_in"`
	RefreshExpiresIn int       `json:"refresh_expires_in"`
	User             UserModel `json:"user"`
}

type LoginReq struct {
	Token    string `json:"token"`
	Provider string `json:"provider"`
}
