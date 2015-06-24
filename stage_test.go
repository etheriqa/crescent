package crescent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedStage struct{ mock.Mock }

func (m *MockedStage) Initialize(g Game) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockedStage) OnTick(g Game) {
	m.Called(g)
}

func TestStageFactories(t *testing.T) {
	assert := assert.New(t)
	s1 := new(MockedStage)
	s2 := new(MockedStage)
	sf := StageFactories{
		1: func() Stage { return s1 },
		2: func() Stage { return s2 },
	}
	assert.Implements((*StageFactory)(nil), sf)
	assert.Equal(s1, sf.New(1))
	assert.Equal(s2, sf.New(2))
	assert.Nil(sf.New(3))
}
