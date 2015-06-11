package main

const (
	DefaultArmor                Statistic = 25
	DefaultMagicResistance      Statistic = 0
	DefaultCriticalStrikeChance Statistic = 0.05
	DefaultCriticalStrikeFactor Statistic = 0.5
	DefaultCooldownReduction    Statistic = 0.0
	DefaultDamageThreatFactor   Statistic = 1.0
	DefaultHealingThreatFactor  Statistic = 0.4
)

type class struct {
	name                 string
	health               Statistic
	healthRegeneration   Statistic
	mana                 Statistic
	manaRegeneration     Statistic
	armor                Statistic
	magicResistance      Statistic
	criticalStrikeChance Statistic
	criticalStrikeFactor Statistic
	cooldownReduction    Statistic
	damageThreatFactor   Statistic
	healingThreatFactor  Statistic
	abilities            []*ability
	initializer          func(u *Unit)
}
