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

// OnAttach removes duplicate Disables
func (h *Disable) OnAttach() {
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
		return
	}

	h.writeOutputUnitAttach()
	h.Object().Register(h)
	h.Object().Dispatch(EventDisabled{})
}

// OnDetach does nothing
func (h *Disable) OnDetach() {
	h.Object().Unregister(h)
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
