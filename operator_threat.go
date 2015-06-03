package main

type threat struct {
	target *unit
	v      int32
}

// isComplete returns always false
func (t *threat) isComplete(u *unit) bool { return false }

// onAttach merges threat operators they have same target
func (t *threat) onAttach(u *unit) {
	for o := range u.operators {
		if o == t {
			continue
		}
		if _, ok := o.(*threat); !ok {
			continue
		}
		if o.(*threat).target != t.target {
			continue
		}
		t.v += o.(*threat).v
		u.detachOperator(o)
	}
}

// onTick does nothing
func (t *threat) onTick(u *unit) {}

// onComplete never triggered
func (t *threat) onComplete(u *unit) {}

// onDetach never triggered
func (t *threat) onDetach(u *unit) {}
