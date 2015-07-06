package game

import (
	"errors"
	"math/rand"
)

type World struct {
	rand  *rand.Rand
	clock InstanceClock
	w     InstanceOutputWriter

	stage   Stage
	effects EffectSet
	units   *UnitMap
}

// NewWorld returns a World
func NewWorld(rand *rand.Rand, clock InstanceClock, w InstanceOutputWriter, stage Stage) *World {
	g := &World{
		rand:  rand,
		clock: clock,
		w:     w,

		stage:   stage,
		effects: MakeEffectSet(),
		units:   NewUnitMap(),
	}
	stage.Initialize(g)
	return g
}

// SyncWorld sends the game state
func (w *World) SyncWorld(writer InstanceOutputWriter) {
	w.units.Each(func(u *Unit) {
		writer.Write(OutputUnitJoin{
			UnitID:    u.ID(),
			UnitGroup: u.Group(),
			UnitName:  u.Name(),
			ClassName: u.ClassName(),
			Health:    u.Health(),
			HealthMax: u.HealthMax(),
			Mana:      u.Mana(),
			ManaMax:   u.ManaMax(),
		})
	})
}

// SyncUnit sends the unit information
func (w *World) SyncUnit(writer InstanceOutputWriter, id UnitID) {
	// TODO refactor
	u := w.units.Find(id)
	// TODO handle the error
	if u == nil {
		return
	}
	as := make([]OutputPlayerAbility, 4)
	for i := 0; i < 4; i++ {
		dts := make([]string, 0)
		for _, dt := range u.Abilities()[i].DisableTypes {
			switch dt {
			case DisableTypeSilence:
				dts = append(dts, "Silence")
			case DisableTypeStun:
				dts = append(dts, "Stun")
			}
		}
		as[i] = OutputPlayerAbility{
			Name:               u.Abilities()[i].Name,
			Description:        u.Abilities()[i].Description,
			TargetType:         u.Abilities()[i].TargetType,
			HealthCost:         u.Abilities()[i].HealthCost,
			ManaCost:           u.Abilities()[i].ManaCost,
			ActivationDuration: u.Abilities()[i].ActivationDuration,
			CooldownDuration:   u.Abilities()[i].CooldownDuration,
			DisableTypes:       dts,
		}
	}
	writer.Write(OutputPlayer{
		UnitID: id,
		Q:      as[0],
		W:      as[1],
		E:      as[2],
		R:      as[3],
	})
}

// Ability activates the ability
func (w *World) Ability(sid UnitID, oid *UnitID, abilityName string) error {
	s := w.units.Find(sid)
	if s == nil {
		return errors.New("Unknown subject UnitID")
	}
	var o *Unit
	if oid != nil {
		o = w.units.Find(*oid)
		if o == nil {
			return errors.New("Unknown object UnitID")
		}
	}
	// TODO refactor
	switch abilityName {
	case "Q":
		w.Activating(s, o, s.Abilities()[0])
	case "W":
		w.Activating(s, o, s.Abilities()[1])
	case "E":
		w.Activating(s, o, s.Abilities()[2])
	case "R":
		w.Activating(s, o, s.Abilities()[3])
	default:
		return errors.New("Unknown ability name")
	}
	return nil
}

// Interrupt interrupts ability activation
func (w *World) Interrupt(id UnitID) error {
	u := w.units.Find(id)
	if u == nil {
		return errors.New("Unknown UnitID")
	}
	u.Dispatch(EventInterrupt{
		UnitID: id,
	})
	return nil
}

// PerformGameTick performs the game tick routine
func (w *World) PerformGameTick() {
	w.stage.OnTick(w)
	w.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventGameTick{})
	})
}

// PerformPeriodicalTick performs the periodical rick routine
func (w *World) PerformPeriodicalTick() {
	w.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.Dispatch(EventPeriodicalTick{})
	})
}

// PerformRegenerationTick performs the regeneration tick routine
func (w *World) PerformRegenerationTick() {
	w.units.Each(func(u *Unit) {
		if u.IsDead() {
			return
		}
		u.ModifyHealth(w.w, u.HealthRegeneration())
		u.ModifyMana(w.w, u.ManaRegeneration())
	})
}

// Rand returns the *rand.Rand
func (w *World) Rand() *rand.Rand {
	return w.rand
}

// Clock returns the InstanceClock
func (w *World) Clock() InstanceClock {
	return w.clock
}

// Writer returns the InstanceOutputWriter
func (w *World) Writer() InstanceOutputWriter {
	return w.w
}

// Join creates a Unit and adds it to the game
func (w *World) Join(group UnitGroup, name UnitName, class *Class) (id UnitID, err error) {
	u, err := w.units.Join(group, name, class)
	if err != nil {
		return
	}
	id = u.ID()
	w.w.Write(OutputUnitJoin{
		UnitID:    u.ID(),
		UnitGroup: u.Group(),
		UnitName:  u.Name(),
		ClassName: u.ClassName(),
		Health:    u.Health(),
		HealthMax: u.HealthMax(),
		Mana:      u.Mana(),
		ManaMax:   u.ManaMax(),
	})
	return
}

// Leave removes the Unit
func (w *World) Leave(id UnitID) (err error) {
	if err = w.units.Leave(id); err != nil {
		return
	}
	w.w.Write(OutputUnitLeave{
		UnitID: id,
	})
	return
}

// UnitQuery returns a UnitQueryable
func (w *World) UnitQuery() UnitQueryable {
	return w.units
}

// AddEffect adds the effect
func (w *World) AttachEffect(e Effect) error {
	return w.effects.Attach(w, e)
}

// RemoveEffect removes the effect
func (w *World) DetachEffect(e Effect) error {
	return w.effects.Detach(w, e)
}

// EffectQuery returns a EffectQueryable
func (w *World) EffectQuery() EffectQueryable {
	return w.effects
}

// Activating attaches a Activating Effect
func (w *World) Activating(s Subject, o *Unit, a *Ability) {
	w.AttachEffect(NewActivating(w, s, o, a, w.clock.Add(a.ActivationDuration)))
}

// Cooldown attaches a Cooldown Effect
func (w *World) Cooldown(o Object, a *Ability) {
	w.AttachEffect(NewCooldown(w, o, a, w.clock.Add(a.CooldownDuration)))
}

// ResetCooldown detaches Cooldown effects
func (w *World) ResetCooldown(o Object, a *Ability) {
	w.AttachEffect(NewCooldown(w, o, a, w.clock.Now()))
}

// Correction attaches a Correction Effect
func (w *World) Correction(o Object, c UnitCorrection, name string, l Statistic, d InstanceDuration) {
	w.AttachEffect(NewCorrection(w, o, c, name, l, w.clock.Add(d)))
}

// Disable attaches a Disable Effect
func (w *World) Disable(o Object, dt DisableType, d InstanceDuration) {
	w.AttachEffect(NewDisable(w, o, dt, w.clock.Add(d)))
}

// DamageThreat attaches a Threat Effect
func (w *World) DamageThreat(s Subject, o Object, d Statistic) {
	w.AttachEffect(NewThreat(w, s, o, d*s.Subject().DamageThreatFactor()))
}

// HealingThreat attaches a Threat Effect
func (w *World) HealingThreat(s Subject, o Object, h Statistic) {
	w.AttachEffect(NewThreat(w, s, o, h*s.Subject().HealingThreatFactor()))
}

// DoT attaches a Periodical Effect
func (w *World) DoT(damage Damage, name string, d InstanceDuration) {
	w.AttachEffect(NewPeriodical(w, damage, damage, name, func() { damage.Perform(w) }, w.clock.Add(d)))
}

// HoT attaches a Periodical Effect
func (w *World) HoT(healing Healing, name string, d InstanceDuration) {
	w.AttachEffect(NewPeriodical(w, healing, healing, name, func() { healing.Perform(w) }, w.clock.Add(d)))
}
