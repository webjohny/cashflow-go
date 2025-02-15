package dto

type SetOptionsLobbyRequestDTO struct {
	HandMode    bool   `json:"handMode" binding:"boolean"`
	MeetLink    string `json:"meetLink,omitempty"`
	BannerLink  string `json:"bannerLink,omitempty"`
	BannerImage string `json:"bannerImage,omitempty"`
	Language    string `json:"language,omitempty"`
}
