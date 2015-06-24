package crescent

type ClassFactory interface {
	New(ClassName) *Class
}

type ClassFactories map[ClassName](func() *Class)

// New creates a Class
func (cf ClassFactories) New(name ClassName) *Class {
	if f, ok := cf[name]; !ok {
		return nil
	} else {
		return f()
	}
}
