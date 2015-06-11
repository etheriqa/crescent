package main

type Damage struct {
	UnitPair
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// NewPhysicalDamage returns a damage affected by armor of the object
func NewPhysicalDamage(up UnitPair, baseDamage Statistic) *Damage {
	return NewTrueDamage(up, baseDamage*up.Object().physicalDamageReductionFactor())
}

// NewMagicDamage returns a damage affected by magic resistance of the object
func NewMagicDamage(up UnitPair, baseDamage Statistic) *Damage {
	return NewTrueDamage(up, baseDamage*up.Object().magicDamageReductionFactor())
}

// NewTrueDamage returns a damage that ignores damage reduction
func NewTrueDamage(up UnitPair, baseDamage Statistic) *Damage {
	return &Damage{
		UnitPair:             up,
		amount:               baseDamage,
		criticalStrikeChance: up.Subject().criticalStrikeChance(),
		criticalStrikeFactor: up.Subject().criticalStrikeFactor(),
	}
}

// NewPureDamage returns a damage that ignores both damage reduction and critical strike
func NewPureDamage(up UnitPair, baseDamage Statistic) *Damage {
	return &Damage{
		UnitPair:             up,
		amount:               baseDamage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Perform subtracts amount the damage from the object and attaches a threat handler to the subject and publishes a message
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		d.amount,
		d.criticalStrikeChance,
		d.criticalStrikeFactor,
	)
	after, before, err = d.Object().modifyHealth(-amount)
	if err != nil {
		return
	}
	if d.Subject() != nil {
		threat := NewDamageThreat(MakeUnitPair(d.Subject(), d.Object()), d.amount)
		d.AttachHandler(threat)
	}
	d.Publish(message{
	// TODO pack message
	})
	return
}
