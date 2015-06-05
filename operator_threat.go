package main

type threat struct {
	unit   *unit
	target *unit
	threat int32
}

// onAttach merges threat operators they have same target
func (t *threat) onAttach() {
	for o := range t.unit.operators {
		if o == t {
			continue
		}
		if _, ok := o.(*threat); !ok {
			continue
		}
		if o.(*threat).target != t.target {
			continue
		}
		t.threat += o.(*threat).threat
		t.unit.detachOperator(o)
	}
}

// onDetach does nothing
func (t *threat) onDetach() {}
