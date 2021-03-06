package game

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
}

type InputDisconnect struct {
}

type InputProfile struct {
	UserName UserName
}

type InputChat struct {
	Message string
}

type InputLevel struct {
	LevelID LevelID
}

type InputJoin struct {
	ClassName ClassName
}

type InputLeave struct {
}

type InputAbility struct {
	AbilityName  string
	ObjectUnitID *UnitID
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
	case "Profile":
		i := InputProfile{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Chat":
		i := InputChat{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Level":
		i := InputLevel{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Join":
		i := InputJoin{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Leave":
		i := InputLeave{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Ability":
		i := InputAbility{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	case "Interrupt":
		i := InputInterrupt{}
		if err = json.Unmarshal(*f.Data, &i); err == nil {
			return i, nil
		}
	}
	return nil, err
}
