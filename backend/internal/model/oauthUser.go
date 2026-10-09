package model

type OauthUser struct {
	ID        int    `json:"id"`
	Username  string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}
