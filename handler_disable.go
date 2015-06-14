package main

type DisableType uint8

const (
	_ DisableType = iota

	DisableTypeSilence
	DisableTypeStun
)

type Disable struct {
	UnitObject
	disableType    DisableType
	expirationTime InstanceTime

	op Operator
}

// OnAttach removes duplicate Disables
func (h *Disable) OnAttach() {
	ok := h.op.Handlers().BindObject(h).Every(func(o Handler) bool {
		switch o := o.(type) {
		case *Disable:
			if h == o || h.disableType != o.disableType {
				return true
			}
			if h.expirationTime <= o.expirationTime {
				return false
			}
			h.op.Handlers().Detach(o)
		}
		return true
	})
	if !ok {
		h.op.Handlers().Detach(h)
		return
	}

	h.writeOutputUnitAttach()
	h.Object().AddEventHandler(h, EventGameTick)
	h.Object().AddEventHandler(h, EventDead)
	h.Object().TriggerEvent(EventDisabled)
}

// OnDetach does nothing
func (h *Disable) OnDetach() {
	h.Object().RemoveEventHandler(h, EventGameTick)
	h.Object().RemoveEventHandler(h, EventDead)
}

// HandleEvent handles the Event
func (h *Disable) HandleEvent(e Event) {
	switch e {
	case EventGameTick:
		if h.op.Clock().Before(h.expirationTime) {
			return
		}
		h.writeOutputUnitDetach()
		h.op.Handlers().Detach(h)
	case EventDead:
		h.op.Handlers().Detach(h)
	}
}

// writeOutputUnitAttach writes a OutputUnitAttach
func (h *Disable) writeOutputUnitAttach() {
	h.op.Writer().Write(OutputUnitAttach{
		UnitID:         h.Object().ID(),
		AttachmentName: "Disable", // TODO
		ExpirationTime: h.expirationTime,
	})
}

// writeOutputUnitDetach writes a OutputUnitDetach
func (h *Disable) writeOutputUnitDetach() {
	h.op.Writer().Write(OutputUnitDetach{
		UnitID:         h.Object().ID(),
		AttachmentName: "Disable", // TODO
	})
}
