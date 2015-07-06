package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitSubject(t *testing.T) {
	assert := assert.New(t)
	u := new(Unit)
	s := MakeSubject(u)
	assert.Implements((*Subject)(nil), s)
	assert.Equal(u, s.Subject())
}

func TestUnitObject(t *testing.T) {
	assert := assert.New(t)
	u := new(Unit)
	o := MakeObject(u)
	assert.Implements((*Object)(nil), o)
	assert.Equal(u, o.Object())
}

func TestPair(t *testing.T) {
	assert := assert.New(t)
	s := new(Unit)
	o := new(Unit)
	p := MakePair(s, o)
	assert.Implements((*Subject)(nil), p)
	assert.Implements((*Object)(nil), p)
	assert.Equal(s, p.Subject())
	assert.Equal(o, p.Object())
}
