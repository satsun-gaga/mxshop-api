package response

import "time"

type UserResponse struct {
	Id int32 `json:"id"`

	NickName string `json:"name"`

	BirthDay time.Time `json:"birthday"`

	Gender string `json:"gender"`

	Mobile string `json:"mobile"`
}
