package main

type ability struct {
	perform      func()
	disableTypes []disableType
	cost         statistic
}

// satisfiedRequirements returns true iff the ability satisfy activation requirements
func (p *ability) satisfiedRequirements(performer *unit) bool {
	if performer.mana() < p.cost {
		return false
	}
	for o := range performer.operators {
		switch o := o.(type) {
		case *cooldown:
			if p == o.ability {
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
	return true
}
