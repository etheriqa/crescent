package main

import (
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

	assert.NotEqual(set, set.Bind(new(Unit)))
	assert.Equal(set, set.Bind(new(Unit)).Unbind())
	assert.Equal(set, set.Bind(new(Unit)).Unbind().Unbind())

	h1 := new(MockedEffect)
	h2 := new(MockedEffect)

	h1.On("EffectDidAttach").Return(nil).Once()
	set.Detach(h1)
	set.Attach(h1)
	set.Attach(h1)
	h1.AssertExpectations(t)

	h2.On("EffectDidAttach").Return(nil).Once()
	set.Detach(h2)
	set.Attach(h2)
	set.Attach(h2)
	h2.AssertExpectations(t)

	h1.On("EffectDidDetach").Return(nil).Once()
	set.Attach(h1)
	set.Detach(h1)
	set.Detach(h1)
	h1.AssertExpectations(t)

	h1.On("EffectDidAttach").Return(nil).Once()
	set.Detach(h1)
	set.Attach(h1)
	set.Attach(h1)
	h1.AssertExpectations(t)
}

func TestEffectSetEach(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	h1 := new(MockedEffect)
	h1.On("EffectDidAttach").Return(nil)
	h2 := new(MockedEffect)
	h2.On("EffectDidAttach").Return(nil)

	{
		hs := []Effect{}
		set.Each(func(h Effect) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		set.Attach(h1)
		hs := []Effect{}
		set.Each(func(h Effect) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, h1)
	}

	{
		set.Attach(h2)
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
	h1.On("EffectDidAttach").Return(nil)
	h2 := new(MockedEffect)
	h2.On("EffectDidAttach").Return(nil)

	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.True(set.Every(func(h Effect) bool { return false }))

	set.Attach(h1)
	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.False(set.Every(func(h Effect) bool { return false }))

	set.Attach(h2)
	assert.True(set.Every(func(h Effect) bool { return true }))
	assert.False(set.Every(func(h Effect) bool { return h.(*MockedEffect) == h1 }))
	assert.False(set.Every(func(h Effect) bool { return h.(*MockedEffect) == h2 }))
	assert.False(set.Every(func(h Effect) bool { return false }))
}

func TestEffectSetSome(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	h1 := new(MockedEffect)
	h1.On("EffectDidAttach").Return(nil)
	h2 := new(MockedEffect)
	h2.On("EffectDidAttach").Return(nil)

	assert.False(set.Some(func(h Effect) bool { return true }))
	assert.False(set.Some(func(h Effect) bool { return false }))

	set.Attach(h1)
	assert.True(set.Some(func(h Effect) bool { return true }))
	assert.False(set.Some(func(h Effect) bool { return false }))

	set.Attach(h2)
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
	h.On("EffectDidAttach").Return(nil).Once()
	sh := new(MockedEffectS)
	sh.On("EffectDidAttach").Return(nil).Once()
	oh := new(MockedEffectO)
	oh.On("EffectDidAttach").Return(nil).Once()
	soh := new(MockedEffectSO)
	soh.On("EffectDidAttach").Return(nil).Once()

	set.Bind(new(Unit)).Attach(h)
	set.Bind(new(Unit)).Attach(sh)
	set.Bind(new(Unit)).Attach(oh)
	set.Bind(new(Unit)).Attach(soh)
	h.AssertExpectations(t)
	sh.AssertExpectations(t)
	oh.AssertExpectations(t)
	soh.AssertExpectations(t)

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

	h.On("EffectDidDetach").Return(nil).Once()
	set.Bind(new(Unit)).Detach(h)
	h.AssertExpectations(t)
}

func TestEffectSetBoundEvery(t *testing.T) {
	assert := assert.New(t)
	set := MakeEffectSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedEffect)
	h.On("EffectDidAttach").Return(nil)
	sh := new(MockedEffectS)
	sh.On("EffectDidAttach").Return(nil)
	sh.On("Subject").Return(s)
	oh := new(MockedEffectO)
	oh.On("EffectDidAttach").Return(nil)
	oh.On("Object").Return(o)
	soh := new(MockedEffectSO)
	soh.On("EffectDidAttach").Return(nil)
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)

	set.Attach(h)
	set.Attach(sh)
	set.Attach(oh)
	set.Attach(soh)

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
	h.On("EffectDidAttach").Return(nil)
	sh := new(MockedEffectS)
	sh.On("EffectDidAttach").Return(nil)
	sh.On("Subject").Return(s)
	oh := new(MockedEffectO)
	oh.On("EffectDidAttach").Return(nil)
	oh.On("Object").Return(o)
	soh := new(MockedEffectSO)
	soh.On("EffectDidAttach").Return(nil)
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)

	set.Attach(h)
	set.Attach(sh)
	set.Attach(oh)
	set.Attach(soh)

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
