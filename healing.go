package main

type Healing struct {
	UnitPair
	healing              Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	g Game
}

// Perform performs the Healing
func (h *Healing) Perform() (before, after Statistic, crit bool, err error) {
	healing, crit := applyCriticalStrike(h.healing, h.criticalStrikeChance, h.criticalStrikeFactor)
	before, after, err = h.Object().ModifyHealth(h.g.Writer(), healing)
	if err != nil {
		return
	}

	h.g.Writer().Write(OutputHealing{
		SubjectUnitID: h.Subject().ID(),
		ObjectUnitID:  h.Object().ID(),
		Healing:       healing,
		IsCritical:    crit,
	})

	h.g.Units().EachEnemy(h.Subject(), func(enemy *Unit) {
		h.g.HealingThreat(h, enemy, healing)
	})
	return
}
