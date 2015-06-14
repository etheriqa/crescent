package main

import (
	"encoding/json"
	"errors"
	"time"
)

type Output struct {
	ClientID ClientID
	Output   interface{}
}

type OutputFrame struct {
	Type string
	Data interface{}
}

type OutputHello struct {
	RealGameTick   time.Duration
	InstanceTime   InstanceTime
	UnitList       []OutputUnitJoin
	AttachmentList []OutputUnitAttach
}

type OutputSync struct {
	InstanceTime InstanceTime
}

type OutputMessage struct {
	Message string
}

type OutputChat struct {
	ClientName ClientName
	Message    string
}

type OutputUnitJoin struct {
	UnitID       UnitID
	UnitGroup    UnitGroup
	UnitPosition UnitPosition
	UnitName     UnitName
	ClassName    string
	Health       Statistic
	HealthMax    Statistic
	Mana         Statistic
	ManaMax      Statistic
}

type OutputUnitLeave struct {
	UnitID
}

type OutputUnitAttach struct {
	UnitID         UnitID
	AttachmentName string
	Stack          Statistic
	ExpirationTime InstanceTime
}

type OutputUnitDetach struct {
	UnitID         UnitID
	AttachmentName string
}

type OutputUnitActivating struct {
	UnitID             UnitID
	AbilityName        string
	ActivationDuration InstanceDuration
	ExpirationTime     InstanceTime
}

type OutputUnitActivated struct {
	UnitID      UnitID
	AbilityName string
	OK          bool
}

type OutputUnitCooldown struct {
	UnitID         UnitID
	AbilityName    string
	ExpirationTime InstanceTime
}

type OutputUnitResource struct {
	UnitID UnitID
	Health Statistic
	Mana   Statistic
}

type OutputDamage struct {
	SubjectUnitID UnitID
	ObjectUnitID  UnitID
	Damage        Statistic
	IsCritical    bool
}

type OutputHealing struct {
	SubjectUnitID UnitID
	ObjectUnitID  UnitID
	Healing       Statistic
	IsCritical    bool
}

// EncodeOutputFrame encodes the Output to a output frame
func EncodeOutputFrame(o interface{}) ([]byte, error) {
	var f OutputFrame
	f.Data = o
	switch o.(type) {
	case OutputHello:
		f.Type = "Hello"
	case OutputSync:
		f.Type = "Sync"
	case OutputMessage:
		f.Type = "Message"
	case OutputChat:
		f.Type = "Chat"
	case OutputUnitJoin:
		f.Type = "UnitJoin"
	case OutputUnitLeave:
		f.Type = "UnitLeave"
	case OutputUnitAttach:
		f.Type = "UnitAttach"
	case OutputUnitDetach:
		f.Type = "UnitDetach"
	case OutputUnitActivating:
		f.Type = "UnitActivating"
	case OutputUnitActivated:
		f.Type = "UnitActivated"
	case OutputUnitCooldown:
		f.Type = "UnitCooldown"
	case OutputUnitResource:
		f.Type = "UnitResource"
	case OutputDamage:
		f.Type = "Damage"
	case OutputHealing:
		f.Type = "Healing"
	default:
		return nil, errors.New("Unknown output type")
	}
	return json.Marshal(f)
}