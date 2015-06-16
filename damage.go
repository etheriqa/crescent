package main

type Damage struct {
	UnitPair
	damage               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	op Operator
}

// Perform performs the Damage
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	damage, crit := applyCriticalStrike(d.damage, d.criticalStrikeChance, d.criticalStrikeFactor)
	before, after, err = d.Object().ModifyHealth(d.op.Writer(), -damage)
	if err != nil {
		return
	}

	d.op.Writer().Write(OutputDamage{
		SubjectUnitID: d.Subject().ID(),
		ObjectUnitID:  d.Object().ID(),
		Damage:        damage,
		IsCritical:    crit,
	})

	if d.Object().IsDead() {
		d.Object().TriggerEvent(EventDead)
		return
	} else {
		d.Object().TriggerEvent(EventTakenDamage)
	}

	d.op.DamageThreat(d, d, damage)
	return
}
