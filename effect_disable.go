package main

type DisableType uint8

const (
	_ DisableType = iota

	DisableTypeSilence
	DisableTypeStun
)

func (d DisableType) String() string {
	switch d {
	case DisableTypeSilence:
		return "Silence"
	case DisableTypeStun:
		return "Stun"
	default:
		return ""
	}
}

type Disable struct {
	UnitObject
	disableType    DisableType
	expirationTime InstanceTime

	op Operator
}

// EffectDidAttach removes duplicate Disables
func (h *Disable) EffectDidAttach() error {
	ok := h.op.Effects().BindObject(h).Every(func(o Effect) bool {
		switch o := o.(type) {
		case *Disable:
			if h == o || h.disableType != o.disableType {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.op.Effects().Detach(o)
		}
		return true
	})
	if !ok {
		h.op.Effects().Detach(h)
		return nil
	}

	h.writeOutputUnitAttach()
	h.Object().Register(h)
	h.Object().Dispatch(EventDisabled{})
	return nil
}

// EffectDidDetach does nothing
func (h *Disable) EffectDidDetach() error {
	h.Object().Unregister(h)
	return nil
}

// Handle handles the Event
func (h *Disable) Handle(p interface{}) {
	switch p.(type) {
	case *EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.writeOutputUnitDetach()
		h.op.Effects().Detach(h)
	case *EventDead:
		h.op.Effects().Detach(h)
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (h *Disable) writeOutputUnitAttach() {
	h.op.Writer().Write(OutputUnitAttach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.disableType.String(),
		ExpirationTime: h.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (h *Disable) writeOutputUnitDetach() {
	h.op.Writer().Write(OutputUnitDetach{
		UnitID:         h.Object().ID(),
		AttachmentName: h.disableType.String(),
	})
}
