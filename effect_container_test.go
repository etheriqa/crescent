package crescent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedEffectContainer struct{ mock.Mock }

func (m *MockedEffectContainer) Attach(h Effect) {
	m.Called(h)
}

func (m *MockedEffectContainer) Detach(h Effect) {
	m.Called(h)
}

func (m *MockedEffectContainer) Bind(u *Unit) EffectContainer {
	args := m.Called(u)
	return args.Get(0).(EffectContainer)
}

func (m *MockedEffectContainer) BindSubject(s Subject) EffectContainer {
	args := m.Called(s)
	return args.Get(0).(EffectContainer)
}

func (m *MockedEffectContainer) BindObject(o Object) EffectContainer {
	args := m.Called(o)
	return args.Get(0).(EffectContainer)
}

func (m *MockedEffectContainer) Unbind() EffectContainer {
	args := m.Called()
	return args.Get(0).(EffectContainer)
}

func (m *MockedEffectContainer) Each(f func(Effect)) {
	m.Called(f)
}

func (m *MockedEffectContainer) Every(f func(Effect) bool) bool {
	args := m.Called(f)
	return args.Bool(0)
}

func (m *MockedEffectContainer) Some(f func(Effect) bool) bool {
	args := m.Called(f)
	return args.Bool(0)
}

func TestEffectSet(t *testing.T) {
	assert := assert.New(t)

	set := MakeEffectSet()
	assert.Implements((*EffectContainer)(nil), set)
	assert.Implements((*EffectQueryable)(nil), set)

	assert.NotEqual(set, set.Bind(new(Unit)))
	assert.Equal(set, set.Bind(new(Unit)).Unbind())
	assert.Equal(set, set.Bind(new(Unit)).Unbind().Unbind())

	e := new(MockedFullEffect)
	g := new(MockedGame)

	e.On("EffectWillAttach", g).Return(nil).Once()
	e.On("EffectDidAttach", g).Return(nil).Once()
	assert.Implements((*error)(nil), set.Detach(g, e))
	assert.Nil(set.Attach(g, e))
	assert.Implements((*error)(nil), set.Attach(g, e))
	e.AssertExpectations(t)

	e.On("EffectWillDetach", g).Return(nil).Once()
	e.On("EffectDidDetach", g).Return(nil).Once()
	assert.Implements((*error)(nil), set.Attach(g, e))
	assert.Nil(set.Detach(g, e))
	assert.Implements((*error)(nil), set.Detach(g, e))
	e.AssertExpectations(t)

	e.On("EffectWillAttach", g).Return(errors.New("error")).Once()
	assert.Implements((*error)(nil), set.Attach(g, e))
	assert.False(set.Some(func(f Effect) bool { return e == f }))
	e.AssertExpectations(t)

	e.On("EffectWillAttach", g).Return(nil).Once()
	e.On("EffectDidAttach", g).Return(errors.New("error")).Once()
	assert.Implements((*error)(nil), set.Attach(g, e))
	assert.True(set.Some(func(f Effect) bool { return e == f }))
	e.AssertExpectations(t)

	e.On("EffectWillDetach", g).Return(errors.New("error")).Once()
	assert.Implements((*error)(nil), set.Detach(g, e))
	assert.True(set.Some(func(f Effect) bool { return e == f }))
	e.AssertExpectations(t)

	e.On("EffectWillDetach", g).Return(nil).Once()
	e.On("EffectDidDetach", g).Return(errors.New("error")).Once()
	assert.Implements((*error)(nil), set.Detach(g, e))
	assert.False(set.Some(func(f Effect) bool { return e == f }))
	e.AssertExpectations(t)
}

func TestEffectSetEach(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	h1 := new(MockedEffect)
	h2 := new(MockedEffect)
	g := new(MockedGame)

	{
		hs := []Effect{}
		set.Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		set.Attach(g, h1)
		hs := []Effect{}
		set.Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, h1)
	}

	{
		set.Attach(g, h2)
		hs := []Effect{}
		set.Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, h1)
		assert.Contains(hs, h2)
	}
}

func TestEffectSetEvery(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	h1 := new(MockedEffect)
	h2 := new(MockedEffect)
	g := new(MockedGame)

	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.True(set.Every(func(h Effect) bool { return false }))

	set.Attach(g, h1)
	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.False(set.Every(func(h Effect) bool { return false }))

	set.Attach(g, h2)
	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.False(set.Every(func(h Effect) bool { return h.(*MockedEffect) == h1 }))
	assert.False(set.Every(func(h Effect) bool { return h.(*MockedEffect) == h2 }))
	assert.False(set.Every(func(h Effect) bool { return false }))
}

func TestEffectSetSome(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	h1 := new(MockedEffect)
	h2 := new(MockedEffect)
	g := new(MockedGame)

	assert.False(set.Some(func(h Effect) bool { return true }))
	assert.False(set.Some(func(h Effect) bool { return false }))

	set.Attach(g, h1)
	assert.True(set.Some(func(h Effect) bool { return true }))
	assert.False(set.Some(func(h Effect) bool { return false }))

	set.Attach(g, h2)
	assert.True(set.Some(func(h Effect) bool { return true }))
	assert.True(set.Some(func(h Effect) bool { return h.(*MockedEffect) == h1 }))
	assert.True(set.Some(func(h Effect) bool { return h.(*MockedEffect) == h2 }))
	assert.False(set.Some(func(h Effect) bool { return false }))
}

func TestEffectSetBoundEach(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedEffect)
	sh := new(MockedEffectS)
	oh := new(MockedEffectO)
	soh := new(MockedEffectSO)
	g := new(MockedGame)

	set.Attach(g, h)
	set.Attach(g, sh)
	set.Attach(g, oh)
	set.Attach(g, soh)
	sh.On("Subject").Return(s)
	oh.On("Object").Return(o)
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)

	{
		hs := []Effect{}
		set.BindSubject(s).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, sh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.BindSubject(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.BindSubject(new(Unit)).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.BindObject(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, oh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.BindObject(s).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.BindObject(new(Unit)).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.Bind(s).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, sh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.Bind(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, oh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.Bind(new(Unit)).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.BindSubject(s).BindObject(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.BindObject(s).BindSubject(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Effect{}
		set.BindSubject(s).Bind(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}

	{
		hs := []Effect{}
		set.Bind(s).BindObject(o).Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}
}

func TestEffectSetBoundEvery(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedEffect)
	sh := new(MockedEffectS)
	sh.On("Subject").Return(s)
	oh := new(MockedEffectO)
	oh.On("Object").Return(o)
	soh := new(MockedEffectSO)
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)
	g := new(MockedGame)

	set.Attach(g, h)
	set.Attach(g, sh)
	set.Attach(g, oh)
	set.Attach(g, soh)

	assert.True(set.Bind(new(Unit)).Every(func(h Effect) bool { return true }))
	assert.True(set.Bind(new(Unit)).Every(func(h Effect) bool { return false }))

	assert.False(set.BindSubject(s).Every(func(h Effect) bool {
		assert.Equal(s, h.(Subject).Subject())
		if _, ok := h.(Object); ok {
			return h.(Object).Object() == o
		} else {
			return false
		}
	}))

	assert.False(set.BindObject(o).Every(func(h Effect) bool {
		assert.Equal(o, h.(Object).Object())
		if _, ok := h.(Subject); ok {
			return h.(Subject).Subject() == s
		} else {
			return false
		}
	}))
}

func TestEffectSetBoundSome(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedEffect)
	sh := new(MockedEffectS)
	sh.On("Subject").Return(s)
	oh := new(MockedEffectO)
	oh.On("Object").Return(o)
	soh := new(MockedEffectSO)
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)
	g := new(MockedGame)

	set.Attach(g, h)
	set.Attach(g, sh)
	set.Attach(g, oh)
	set.Attach(g, soh)

	assert.False(set.Bind(new(Unit)).Some(func(h Effect) bool { return true }))
	assert.False(set.Bind(new(Unit)).Some(func(h Effect) bool { return false }))

	assert.True(set.BindSubject(s).Some(func(h Effect) bool {
		assert.Equal(s, h.(Subject).Subject())
		if _, ok := h.(Object); ok {
			return h.(Object).Object() == o
		} else {
			return false
		}
	}))

	assert.True(set.BindObject(o).Some(func(h Effect) bool {
		assert.Equal(o, h.(Object).Object())
		if _, ok := h.(Subject); ok {
			return h.(Subject).Subject() == s
		} else {
			return false
		}
	}))
}
