package main

func newClassHealer() *class {
	class := &class{
		name: "Healer",
		// TODO stats
		health:               700,
		healthRegeneration:   2,
		mana:                 400,
		manaRegeneration:     6,
		armor:                defaultArmor,
		magicResistance:      defaultMagicResistance,
		criticalStrikeChance: defaultCriticalStrikeChance,
		criticalStrikeFactor: defaultCriticalStrikeFactor,
		cooldownReduction:    defaultCooldownReduction,
		damageThreatFactor:   defaultDamageThreatFactor,
		healingThreatFactor:  defaultHealingThreatFactor,
		// TODO abilities
		abilities: make([]ability, 4),
	}
	return class
}
