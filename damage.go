package crescent

type Damage struct {
	UnitPair
	damage               Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
}

// MakePhysicalDamage returns a physical Damage
func MakePhysicalDamage(s Subject, o Object, damage Statistic) Damage {
	return Damage{
		UnitPair:             MakePair(s, o),
		damage:               damage * o.Object().PhysicalDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),
	}
}

// MakeMagicDamage returns a magic Damage
func MakeMagicDamage(s Subject, o Object, damage Statistic) Damage {
	return Damage{
		UnitPair:             MakePair(s, o),
		damage:               damage * o.Object().MagicDamageReductionFactor(),
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),
	}
}

// MakeTrueDamage returns a true Damage
func MakeTrueDamage(s Subject, o Object, damage Statistic) Damage {
	return Damage{
		UnitPair:             MakePair(s, o),
		damage:               damage,
		criticalStrikeChance: s.Subject().CriticalStrikeChance(),
		criticalStrikeFactor: s.Subject().CriticalStrikeFactor(),
	}
}

// MakePureDamage returns a pure Damage
func MakePureDamage(s Subject, o Object, damage Statistic) Damage {
	return Damage{
		UnitPair:             MakePair(s, o),
		damage:               damage,
		criticalStrikeChance: 0,
		criticalStrikeFactor: 0,
	}
}

// Perform performs the Damage
func (d Damage) Perform(g Game) (before, after Statistic, crit bool, err error) {
	damage, crit := applyCriticalStrike(g.Rand(), d.damage, d.criticalStrikeChance, d.criticalStrikeFactor)
	before, after, err = d.Object().ModifyHealth(g.Writer(), -damage)
	if err != nil {
		return
	}

	g.Writer().Write(OutputDamage{
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

	g.DamageThreat(d, d, damage)
	return
}
