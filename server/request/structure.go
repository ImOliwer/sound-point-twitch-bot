package request

type TwitchUser struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       uint64 `json:"view_count"`
	CreatedAt       string `json:"created_at"`
}

type TwitchUserList struct {
	Users []TwitchUser `json:"data"`
}
