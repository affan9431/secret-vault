package models

type Secrets struct {
	Id int64 `json:"id"`
	UserId int64 `json:"user_id"`
	Title string `json:"title"`
	Secret string `json:"secret"`
	Tags string `json:"tags"`
	ExtraData *string `json:"extra_data,omitempty"`
}

