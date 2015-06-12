package main

type Healing struct {
	UnitPair
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// NewHealing returns a healing
func NewHealing(up UnitPair, baseHealing Statistic) *Healing {
	return &Healing{
		UnitPair:             up,
		amount:               baseHealing,
		criticalStrikeChance: up.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: up.Subject().CriticalStrikeFactor(),
	}
}

// NewPureHealing returns a healing that ignores critical strike
func NewPureHealing(up UnitPair, baseHealing Statistic) *Healing {
	return &Healing{
		UnitPair:             up,
		amount:               baseHealing,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Perform adds amount of healing to the object and attaches a threat handler to the enemies and publish a message
func (h *Healing) Perform() (after, before Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		h.amount,
		h.criticalStrikeChance,
		h.criticalStrikeFactor,
	)
	after, before, err = h.Object().ModifyHealth(amount)
	if err != nil {
		return
	}
	if h.Subject() != nil {
		for _, enemy := range h.Subject().Enemies() {
			threat := NewHealingThreat(MakeUnitPair(h.Subject(), enemy), h.amount)
			h.AttachHandler(threat)
		}
	}
	h.Publish(message{
	// TODO pack message
	})
	return
}
