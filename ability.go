package main

type ability interface {
	perform()
	satisfiedRequirements(u *unit) bool
}

type partialAbility struct {
	disableTypes []disableType
}

// perform does nothing
func (p *partialAbility) perform() {}

// satisfiedRequirements returns true iff the ability satisfy activation requirements
func (p *partialAbility) satisfiedRequirements(u *unit) bool {
	for o := range u.operators {
		switch o := o.(type) {
		case *cooldown:
			if _, ok := o.ability.(*partialAbility); !ok {
				continue
			}
			if p == o.ability.(*partialAbility) {
				return false
			}
		case *disable:
			for d := range p.disableTypes {
				if disableType(d) == o.disableType {
					return false
				}
			}
		}
	}
	// todo check cost
	return true
}
