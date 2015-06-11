package main

type Damage struct {
	subject              *Unit
	object               *Unit
	amount               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// NewPhysicalDamage returns a damage affected by armor of the object
func NewPhysicalDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return NewTrueDamage(
		subject,
		object,
		baseDamage*object.physicalDamageReductionFactor(),
	)
}

// NewMagicDamage returns a damage affected by magic resistance of the object
func NewMagicDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return NewTrueDamage(
		subject,
		object,
		baseDamage*object.magicDamageReductionFactor(),
	)
}

// NewTrueDamage returns a damage that ignores damage reduction
func NewTrueDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return &Damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: subject.criticalStrikeChance(),
		criticalStrikeFactor: subject.criticalStrikeFactor(),
	}
}

// NewPureDamage returns a damage that ignores both damage reduction and critical strike
func NewPureDamage(subject, object *Unit, baseDamage Statistic) *Damage {
	return &Damage{
		subject:              subject,
		object:               object,
		amount:               baseDamage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Subject returns the subject
func (d *Damage) Subject() *Unit {
	return d.subject
}

// Object returns the object
func (d *Damage) Object() *Unit {
	return d.object
}

// Perform subtracts amount the damage from the object and attaches a threat handler to the subject and publishes a message
func (d *Damage) Perform() (before, after Statistic, crit bool, err error) {
	amount, crit := applyCriticalStrike(
		d.amount,
		d.criticalStrikeChance,
		d.criticalStrikeFactor,
	)
	after, before, err = d.object.modifyHealth(-amount)
	if err != nil {
		return
	}
	if d.subject != nil {
		d.object.AttachHandler(newDamageThreat(d.subject, d.object, d.amount))
	}
	d.object.Publish(message{
	// TODO pack message
	})
	return
}
