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

type Class struct {
	Name                 string
	Health               Statistic
	HealthRegeneration   Statistic
	Mana                 Statistic
	ManaRegeneration     Statistic
	Armor                Statistic
	MagicResistance      Statistic
	CriticalStrikeChance Statistic
	CriticalStrikeFactor Statistic
	CooldownReduction    Statistic
	DamageThreatFactor   Statistic
	HealingThreatFactor  Statistic
	Abilities            []*Ability
	Initializer          func(u *Unit)
}

func (c *Class) Ability(key string) *Ability {
	switch key {
	case "q":
		return c.Abilities[0]
	case "w":
		return c.Abilities[1]
	case "e":
		return c.Abilities[2]
	case "r":
		return c.Abilities[3]
	default:
		return nil
	}
}
