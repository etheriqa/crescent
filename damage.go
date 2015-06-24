package crescent

type Damage struct {
	UnitPair
	damage               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	g Game
}

// Perform performs the Damage
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	damage, crit := applyCriticalStrike(d.g.Rand(), d.damage, d.criticalStrikeChance, d.criticalStrikeFactor)
	before, after, err = d.Object().ModifyHealth(d.g.Writer(), -damage)
	if err != nil {
		return
	}

	d.g.Writer().Write(OutputDamage{
		SubjectUnitID: d.Subject().ID(),
		ObjectUnitID:  d.Object().ID(),
		Damage:        damage,
		IsCritical:    crit,
	})

	if d.Object().IsDead() {
		d.Object().Dispatch(EventDead{})
		return
	} else {
		d.Object().Dispatch(EventTakenDamage{})
	}

	d.g.DamageThreat(d, d, damage)
	return
}
