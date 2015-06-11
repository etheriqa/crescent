package main

func newClassHealer() *class {
	var q, w, e, r *ability
	class := &class{
		name: "Healer",
		// TODO stats
		health:               700,
		healthRegeneration:   2,
		mana:                 400,
		manaRegeneration:     6,
		armor:                DefaultArmor,
		magicResistance:      DefaultMagicResistance,
		criticalStrikeChance: DefaultCriticalStrikeChance,
		criticalStrikeFactor: DefaultCriticalStrikeFactor,
		cooldownReduction:    DefaultCooldownReduction,
		damageThreatFactor:   DefaultDamageThreatFactor,
		healingThreatFactor:  DefaultHealingThreatFactor,
	}
	// Magic damage / Mana restore
	q = &ability{
		name:               "Healer Q",
		TargetType:         TargetTypeEnemy,
		healthCost:         0,
		manaCost:           0,
		activationDuration: 2 * Second,
		cooldownDuration:   2 * Second,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			// TODO handle the error
			before, after, _, _ := NewMagicDamage(subject, object, 100).Perform()
			// TODO send a message including the ability name
			subject.performManaModification((before - after) * 0.1)
		},
	}
	// HoT
	w = &ability{
		name:               "Healer W",
		TargetType:         TargetTypeFriend,
		healthCost:         0,
		manaCost:           40,
		activationDuration: 2 * Second,
		cooldownDuration:   4 * Second,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			object.AttachHandler(NewHoT(
				NewHealing(subject, object, 20),
				w,
				12*Second,
			))
		},
	}
	// Healing
	e = &ability{
		name:               "Healer E",
		TargetType:         TargetTypeFriend,
		healthCost:         0,
		manaCost:           80,
		activationDuration: 2 * Second,
		cooldownDuration:   8 * Second,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			NewHealing(subject, object, 400).Perform()
		},
	}
	// HoT / Increasing critical strike chance and critical strike factor
	r = &ability{
		name:               "Healer R",
		TargetType:         TargetTypeNone,
		healthCost:         0,
		manaCost:           200,
		activationDuration: 0,
		cooldownDuration:   60 * Second,
		disableTypes: []DisableType{
			DisableTypeSilence,
			DisableTypeStun,
		},
		Perform: func(subject, object *Unit) {
			subject.AttachHandler(NewCorrector(
				subject,
				UnitCorrection{
					CriticalStrikeChance: 0.5,
					CriticalStrikeFactor: 1.5,
				},
				r.name,
				1,
				6*Second,
			))
			for _, friend := range subject.Friends() {
				friend.AttachHandler(NewHoT(
					NewHealing(subject, friend, 20),
					r,
					6*Second,
				))
			}
		},
	}
	class.abilities = []*ability{q, w, e, r}
	return class
}
