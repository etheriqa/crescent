package main

type subject struct {
	observers map[observer]bool
}

func newSubject() *subject {
	return &subject{
		observers: make(map[observer]bool),
	}
}

func (s *subject) attach(o observer) {
	s.observers[o] = true
}

func (s *subject) detach(o observer) {
	delete(s.observers, o)
}

func (s *subject) notify() {
	for o := range s.observers {
		o.update()
	}
}

type observer interface {
	update()
}
