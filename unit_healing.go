package main

type Healing struct {
	UnitPair
	healing              Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	handlers HandlerContainer
	operator Operator
	units    UnitContainer
	writer   GameEventWriter
}

// Perform performs the Healing
func (h *Healing) Perform() (before, after Statistic, crit bool, err error) {
	healing, crit := applyCriticalStrike(h.healing, h.criticalStrikeChance, h.criticalStrikeFactor)
	after, before, err = h.Object().ModifyHealth(h.writer, healing)
	if err != nil {
		log.Fatal(err)
	}

	h.units.EachEnemy(h.Subject(), func(enemy *Unit) {
		h.operator.HealingThreat(h, enemy, healing)
	})

	h.writer.Write(nil) // TODO
	return
}
