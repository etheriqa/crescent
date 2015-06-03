package main

type subject struct {
	observers map[observer]interface{}
}

func newSubject() *subject {
	return &subject{
		observers: make(map[observer]interface{}),
	}
}

func (s *subject) attach(o observer) {
	s.observers[o] = nil
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
