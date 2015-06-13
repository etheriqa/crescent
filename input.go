package main

import (
	"encoding/json"
	"errors"
)

type Input struct {
	ClientID ClientID
	Input    interface{}
}

type InputFrame struct {
	Type string
	Data *json.RawMessage
}

type InputConnect struct {
	ClientName ClientName
}

type InputDisconnect struct {
}

type InputChat struct {
	Message string
}

type InputStage struct {
	StageID StageID
}

type InputJoin struct {
}

type InputLeave struct {
}

type InputAbility struct {
	AblityName   string
	ObjectUnitID UnitID
}

type InputInterrupt struct {
}

// DecodeInputFrame decodes the input frame to a Input
func DecodeInputFrame(p []byte) (interface{}, error) {
	var f InputFrame
	if err := json.Unmarshal(p, &f); err != nil {
		return nil, err
	}
	var i InputChat
	switch f.Type {
	case "Chat":
		i = InputChat{}
		/*
			case "Stage":
				i = InputStage{}
			case "Join":
				i = InputJoin{}
			case "Leave":
				i = InputLeave{}
			case "Ability":
				i = InputAbility{}
			case "Interrupt":
				i = InputInterrupt{}
		*/
	default:
		return nil, errors.New("Unknown input type: " + f.Type)
	}
	if err := json.Unmarshal(*f.Data, &i); err != nil {
		return nil, err
	}
	return i, nil
}
