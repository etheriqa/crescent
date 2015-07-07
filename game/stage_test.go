package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedLevel struct{ mock.Mock }

func (m *MockedLevel) Initialize(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockedLevel) OnTick(g Game) {
	m.Called(g)
}

func TestLevelFactories(t *testing.T) {
	assert := assert.New(t)
	s1 := new(MockedLevel)
	s2 := new(MockedLevel)
	sf := LevelFactories{
		1: func() Level { return s1 },
		2: func() Level { return s2 },
	}
	assert.Implements((*LevelFactory)(nil), sf)
	assert.Equal(s1, sf.New(1))
	assert.Equal(s2, sf.New(2))
	assert.Nil(sf.New(3))
}
