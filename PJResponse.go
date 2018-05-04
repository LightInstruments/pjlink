package pjlink

import (
	"errors"
	"log"
	"strings"
)

type PJResponse struct {
	Class    string   `json:"class"`
	Command  string   `json:"command"`
	Response []string `json:"response"`
}

func NewPJResponse() *PJResponse {
	return &PJResponse{}
}

func (res *PJResponse) Parse(raw string) error {
	// If password is wrong, response will be 'PJLINK ERRA'
	if strings.Contains(raw, "ERRA") { //if authentication succeeded
		return errors.New("Incorrect password")
		//example response: "%1POWR=0"
		//returned params are class, command, and response code(s), respectively
	}
	log.Print(raw)

	tokens := strings.Split(raw, " ")

	token0 := tokens[0]
	param1 := []string{token0[7:len(token0)]}
	paramsN := tokens[1:len(tokens)]
	params := append(param1, paramsN...)

	res.Class = token0[1:2]
	res.Command = token0[2:6]
	res.Response = params

	return nil
}
