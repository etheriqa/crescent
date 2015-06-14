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
	after, before, err = d.Object().ModifyHealth(d.op.Writer(), -damage)
	if err != nil {
		log.Fatal(err)
	}

	d.op.DamageThreat(d, d, damage)
	d.op.Writer().Write(nil) // TODO
	return
}
