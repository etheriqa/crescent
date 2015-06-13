package main

type Damage struct {
	UnitPair
	damage               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic

	handlers HandlerContainer
	operator Operator
	writer   GameEventWriter
}

// Perform performs the Damage
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	damage, crit := applyCriticalStrike(d.damage, d.criticalStrikeChance, d.criticalStrikeFactor)
	after, before, err = d.Object().ModifyHealth(d.writer, -damage)
	if err != nil {
		log.Fatal(err)
	}

	d.operator.DamageThreat(d, d, damage)
	d.writer.Write(nil) // TODO
	return
}
