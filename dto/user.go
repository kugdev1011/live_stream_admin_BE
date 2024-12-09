package dto

type UserQuery struct {
	Role string `query:"role" validate:"omitempty,oneof=supper_admin admin streamer user"`
}
