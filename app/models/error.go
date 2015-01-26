package models

import (
	uuid "github.com/satori/go.uuid"
	"github.com/JohnSmithX/mus/app/utils"
	"encoding/json"
	"time"
)

type ErrMsg struct {
	Id 					uuid.UUID		`json:"id"`
	CreateTime			utils.Time		`json:"create_at"`
	UpdateTime			utils.Time		`json:"update_at"`
	Message				string			`json:"msg"`
}

func (self *ErrMsg) JSON() (result []byte, err error) {
	result, err = json.Marshal(self)
	return
}


func NewErr(msg string) (err *ErrMsg) {
	return &ErrMsg{
		Id: uuid.NewV4(),
		CreateTime: utils.Time(time.Now()),
		UpdateTime: utils.Time(time.Now()),
		Message: msg,
	}
}
