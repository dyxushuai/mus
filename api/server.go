package api

type Server struct {
	ID string `json:"id"` //uuid
	Port string `json:"port" binding:"required"`
	State string `json:"state"`
	Create int64 `json:"create_at"`
	Update int64 `json:"update_at"`

}
