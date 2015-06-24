package crescent

const (
	DefaultArmor                Statistic = 25
	DefaultMagicResistance      Statistic = 0
	DefaultCriticalStrikeChance Statistic = 0.05
	DefaultCriticalStrikeFactor Statistic = 0.5
	DefaultCooldownReduction    Statistic = 0.0
	DefaultDamageThreatFactor   Statistic = 1.0
	DefaultHealingThreatFactor  Statistic = 0.6
)

type ClassName string

type Class struct {
	Name                 ClassName
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
}

type ClassFactory interface {
	New(ClassName) *Class
}

type ClassFactories map[ClassName](func() *Class)

// Ability returns the Ability
func (c *Class) Ability(name string) *Ability {
	for _, a := range c.Abilities {
		if a.Name == name {
			return a
		}
	}
	return nil
}

// New creates a Class
func (cf ClassFactories) New(name ClassName) *Class {
	if f, ok := cf[name]; !ok {
		return nil
	} else {
		return f()
	}
}
