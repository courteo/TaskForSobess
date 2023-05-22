package forms

import (
	"time"
)

type SendNameForm struct {
	Name string 		`json:"name"`
}
type SendTimeForm struct {
	Time time.Duration 	`json:"time"`
}

type LoginForm struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}