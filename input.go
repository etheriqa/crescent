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
	ClassName string
}

type InputLeave struct {
}

type InputAbility struct {
	AbilityName  string
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
	if f.Data == nil {
		return nil, errors.New("Input frame does not have data key")
	}
	err := errors.New("Unknown input type: " + f.Type)
	switch f.Type {
	case "Chat":
		i := InputChat{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Stage":
		i := InputStage{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Join":
		i := InputJoin{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Leave":
		i := InputLeave{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Ability":
		i := InputAbility{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Interrupt":
		i := InputInterrupt{}
		if err := json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	}
	return nil, err
}
