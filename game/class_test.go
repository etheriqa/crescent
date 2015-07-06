package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassFactories(t *testing.T) {
	assert := assert.New(t)
	c1 := new(Class)
	c2 := new(Class)
	cf := ClassFactories{
		"Class1": func() *Class { return c1 },
		"Class2": func() *Class { return c2 },
	}
	assert.Implements((*ClassFactory)(nil), cf)
	assert.Equal(c1, cf.New("Class1"))
	assert.Equal(c2, cf.New("Class2"))
	assert.Nil(cf.New("Class3"))
}
