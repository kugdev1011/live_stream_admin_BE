package dto

type LiveStreamRespDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	VideoSize   string `json:"video_size,omitempty"`
	Likes       uint   `json:"likes,omitempty"`
	Viewers     uint   `json:"viewers,omitempty"`
	Comments    uint   `json:"comments,omitempty"`
	Duration    string `json:"duration,omitempty"`
}
