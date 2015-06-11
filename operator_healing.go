package main

type Healing struct {
	subject              *Unit
	object               *Unit
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// NewHealing returns a healing
func NewHealing(subject, object *Unit, baseHealing Statistic) *Healing {
	return &Healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// NewPureHealing returns a healing that ignores critical strike
func NewPureHealing(subject, object *Unit, baseHealing Statistic) *Healing {
	return &Healing{
		subject:              subject,
		object:               object,
		amount:               baseHealing,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Subject returns the subject
func (h *Healing) Subject() *Unit {
	return h.subject
}

// Object returns the object
func (h *Healing) Object() *Unit {
	return h.object
}

// Perform adds amount of healing to the object and attaches a threat handler to the enemies and publish a message
func (h *Healing) Perform() (after, before Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		h.amount,
		h.criticalStrikeChance,
		h.criticalStrikeFactor,
	)
	after, before, err = h.object.modifyHealth(amount)
	if err != nil {
		return
	}
	if h.subject != nil {
		for _, enemy := range h.subject.Enemies() {
			enemy.AttachHandler(NewHealingThreat(h.subject, enemy, h.amount))
		}
	}
	h.object.Publish(message{
	// TODO pack message
	})
	return
}
