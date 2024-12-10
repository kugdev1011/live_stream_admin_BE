package dto

type UserQuery struct {
	Role string `query:"role" validate:"omitempty,oneof=super_admin admin streamer user"`
}
