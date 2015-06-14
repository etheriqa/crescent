package main

type Healing struct {
	UnitPair
	healing              Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	op Operator
}

// Perform performs the Healing
func (h *Healing) Perform() (before, after Statistic, crit bool, err error) {
	healing, crit := applyCriticalStrike(h.healing, h.criticalStrikeChance, h.criticalStrikeFactor)
	after, before, err = h.Object().ModifyHealth(h.op.Writer(), healing)
	if err != nil {
		log.Fatal(err)
	}

	h.op.Units().EachEnemy(h.Subject(), func(enemy *Unit) {
		h.op.HealingThreat(h, enemy, healing)
	})

	h.op.Writer().Write(OutputHealing{
		SubjectUnitID: h.Subject().ID(),
		ObjectUnitID:  h.Object().ID(),
		Healing:       healing,
		IsCritical:    crit,
	})
	return
}
