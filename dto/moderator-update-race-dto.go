package dto

type ModeratorUpdateRaceDto struct {
	Status        string `json:"status" binding:"required"`
	HideCards     bool   `json:"hide_cards" binding:"boolean"`
	HandMode      bool   `json:"hand_mode" binding:"boolean"`
	MeetLink      string `json:"meet_link,omitempty"`
	BannerLink    string `json:"banner_link,omitempty"`
	BannerImage   string `json:"banner_image,omitempty"`
	EnableManager bool   `json:"enable_manager" binding:"boolean"`
	CurrentPlayer int    `json:"current_player" binding:"numeric"`
	Responses     []bool `json:"responses" binding:"required"`
}
