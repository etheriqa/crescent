package crescent

type Healing struct {
	UnitPair
	healing              Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// NewHealing returns a Healing
func NewHealing(s Subject, o Object, healing Statistic) *Healing {
	return &Healing{
		UnitPair:             MakePair(s, o),
		healing:              healing,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),
	}
}

// Perform performs the Healing
func (h *Healing) Perform(g Game) (before, after Statistic, crit bool, err error) {
	healing, crit := applyCriticalStrike(g.Rand(), h.healing, h.criticalStrikeChance, h.criticalStrikeFactor)
	before, after, err = h.Object().ModifyHealth(g.Writer(), healing)
	if err != nil {
		return
	}

	g.Writer().Write(OutputHealing{
		SubjectUnitID: h.Subject().ID(),
		ObjectUnitID:  h.Object().ID(),
		Healing:       healing,
		IsCritical:    crit,
	})

	g.UnitQuery().EachEnemy(h.Subject(), func(enemy *Unit) {
		g.HealingThreat(h, enemy, healing)
	})
	return
}
