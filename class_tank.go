package main

func newClassTank() *class {
	class := &class{
		name: "Tank",
		// TODO stats
		health:               1000,
		healthRegeneration:   5,
		mana:                 200,
		manaRegeneration:     3,
		armor:                defaultArmor + 25,
		magicResistance:      defaultMagicResistance + 25,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor + 2.0,
		healingThreatFactor:  defaultHealingThreatFactor,
		// TODO abilities
		abilities: make([]*ability, 4),
	}
	return class
}
