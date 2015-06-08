package main

const (
	defaultArmor                statistic = 25
	defaultMagicResistance      statistic = 0
	defaultCriticalStrikeChance statistic = 0.05
	defaultCriticalStrikeFactor statistic = 0.5
	defaultCooldownReduction    statistic = 0.0
	defaultDamageThreatFactor   statistic = 1.0
	defaultHealingThreatFactor  statistic = 0.4
)

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
	abilities            []*ability
	initializer          func(u *unit)
}
