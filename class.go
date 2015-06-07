package main

type class struct {
	name                 string
	health               statistic
	healthRegeneration   statistic
	mana                 statistic
	manaRegeneration     statistic
	armor                statistic
	magicResistance      statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
	cooldownReduction    statistic
	damageThreatFactor   statistic
	healingThreatFactor  statistic
	abilities            []ability
}
