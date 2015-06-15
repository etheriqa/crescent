package main

import (
	"errors"
)

type UnitID uint64
type UnitName string
type UnitGroup uint8
type UnitPosition uint8

type Subject interface {
	Subject() *Unit
}

type Object interface {
	Object() *Unit
}

type Unit struct {
	id         UnitID
	group      UnitGroup
	position   UnitPosition
	name       UnitName
	resource   UnitResource
	correction UnitCorrection
	class      *Class

	EventDispatcher
}

type UnitResource struct {
	Health Statistic
	Mana   Statistic
}

type UnitCorrection struct {
	Armor                Statistic
	MagicResistance      Statistic
	CriticalStrikeChance Statistic
	CriticalStrikeFactor Statistic
	CooldownReduction    Statistic
	DamageThreatFactor   Statistic
	HealingThreatFactor  Statistic
}

// NewUnit returns a Unit
func NewUnit(id UnitID, group UnitGroup, position UnitPosition, name UnitName, class *Class) *Unit {

	return &Unit{
		id:         id,
		name:       name,
		group:      group,
		position:   position,
		resource:   MakeUnitResource(class),
		correction: MakeUnitCorrection(),
		class:      class,

		EventDispatcher: MakeEventHandlerSet(),
	}
}

// Subject returns self
func (u *Unit) Subject() *Unit {
	return u
}

// Object returns self
func (u *Unit) Object() *Unit {
	return u
}

// ID returns the UnitID
func (u *Unit) ID() UnitID {
	return u.id
}

// Group returns the UnitGroup
func (u *Unit) Group() UnitGroup {
	return u.group
}

// Position returns the UnitPosition
func (u *Unit) Position() UnitPosition {
	return u.position
}

// Name returns the UnitName
func (u *Unit) Name() UnitName {
	return u.name
}

// ClassName returns the ClassName
func (u *Unit) ClassName() ClassName {
	return u.class.Name
}

// IsAlive returns true if the Unit is alive
func (u *Unit) IsAlive() bool {
	return u.resource.Health > 0
}

// IsDead returns true if the Unit is Dead
func (u *Unit) IsDead() bool {
	return u.resource.Health <= 0
}

// Health returns amount of health
func (u *Unit) Health() Statistic {
	return u.resource.Health
}

// HealthMax returns maximum amount of health
func (u *Unit) HealthMax() Statistic {
	return u.class.Health
}

// HealthRegeneration returns regeneration speed of health
func (u *Unit) HealthRegeneration() Statistic {
	return u.class.HealthRegeneration
}

// Mana returns amount of mana
func (u *Unit) Mana() Statistic {
	return u.resource.Mana
}

// ManaMax returns maximum amount of mana
func (u *Unit) ManaMax() Statistic {
	return u.class.Mana
}

// ManaRegeneration return regeneration speed of mana
func (u *Unit) ManaRegeneration() Statistic {
	return u.class.ManaRegeneration
}

// Armor returns amount of armor
func (u *Unit) Armor() Statistic {
	return u.class.Armor + u.correction.Armor
}

// MagicResistance returns amount of MagicResistance
func (u *Unit) MagicResistance() Statistic {
	return u.class.MagicResistance + u.correction.MagicResistance
}

// PhysicalDamageReductionFactor returns damage reduction factor for physical damage
func (u *Unit) PhysicalDamageReductionFactor() Statistic {
	return damageReductionFactor(u.Armor())
}

// MagicDamageReductionFactor returns damage reduction factor for magic damage
func (u *Unit) MagicDamageReductionFactor() Statistic {
	return damageReductionFactor(u.MagicResistance())
}

// CriticalStrikeChance returns critical strike chance
func (u *Unit) CriticalStrikeChance() Statistic {
	return u.class.CriticalStrikeChance + u.correction.CriticalStrikeChance
}

// CriticalStrikeFactor returns critical strike factor
func (u *Unit) CriticalStrikeFactor() Statistic {
	return u.class.CriticalStrikeFactor + u.correction.CriticalStrikeFactor
}

// CooldownReduction returns amount of cooldown reduction
func (u *Unit) CooldownReduction() Statistic {
	return u.class.CooldownReduction + u.correction.CooldownReduction
}

// DamageThreatFactor returns threat factor for dealing damage
func (u *Unit) DamageThreatFactor() Statistic {
	return u.class.DamageThreatFactor + u.correction.DamageThreatFactor
}

// HealingThreatFactor returns threat factor for dealing healing
func (u *Unit) HealingThreatFactor() Statistic {
	return u.class.HealingThreatFactor + u.correction.HealingThreatFactor
}

// Ability returns the ability
func (u *Unit) Ability(name string) *Ability {
	return u.class.Ability(name)
}

// UpdateCorrection updates the UnitCorrection
func (u *Unit) UpdateCorrection(correction UnitCorrection) {
	u.correction = correction
}

// ModifyHealth modifies health and returns before/after amount of health
func (u *Unit) ModifyHealth(w InstanceOutputWriter, delta Statistic) (before, after Statistic, err error) {
	if u.IsDead() {
		err = errors.New("Cannot modify health of dead unit")
		return
	}
	before = u.Health()
	after = u.Health() + delta
	if after < 0 {
		after = 0
	}
	if after > u.HealthMax() {
		after = u.HealthMax()
	}
	if before == after {
		return
	}
	u.resource.Health = after
	u.writeOutputUnitResource(w)
	return
}

// ModifyMana modifies mana and returns before/after amount of mana
func (u *Unit) ModifyMana(w InstanceOutputWriter, delta Statistic) (before, after Statistic, err error) {
	if u.IsDead() {
		err = errors.New("Cannot modify mana of dead unit")
		return
	}
	before = u.Mana()
	after = u.Mana() + delta
	if after < 0 {
		after = 0
	}
	if after > u.ManaMax() {
		after = u.ManaMax()
	}
	if before == after {
		return
	}
	u.resource.Mana = after
	u.writeOutputUnitResource(w)
	return
}

// writeOutputUnitResource write a OutputUnitResource
func (u *Unit) writeOutputUnitResource(w InstanceOutputWriter) {
	w.Write(OutputUnitResource{
		UnitID: u.id,
		Health: u.Health(),
		Mana:   u.Mana(),
	})
}

// MakeUnitResource returns a UnitResource
func MakeUnitResource(class *Class) UnitResource {
	return UnitResource{
		Health: class.Health,
		Mana:   class.Mana,
	}
}

// MakeUnitCorrection returns a UnitCorrection
func MakeUnitCorrection() UnitCorrection {
	return UnitCorrection{}
}
