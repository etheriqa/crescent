package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerSet(t *testing.T) {
	assert := assert.New(t)

	set := MakeHandlerSet()
	assert.Implements((*HandlerContainer)(nil), set)

	assert.NotEqual(set, set.Bind(new(Unit)))
	assert.Equal(set, set.Bind(new(Unit)).Unbind())
	assert.Equal(set, set.Bind(new(Unit)).Unbind().Unbind())

	h1 := new(MockedHandler)
	h2 := new(MockedHandler)

	h1.On("OnAttach").Return().Once()
	set.Detach(h1)
	set.Attach(h1)
	set.Attach(h1)
	h1.AssertExpectations(t)

	h2.On("OnAttach").Return().Once()
	set.Detach(h2)
	set.Attach(h2)
	set.Attach(h2)
	h2.AssertExpectations(t)

	h1.On("OnDetach").Return().Once()
	set.Attach(h1)
	set.Detach(h1)
	set.Detach(h1)
	h1.AssertExpectations(t)

	h1.On("OnAttach").Return().Once()
	set.Detach(h1)
	set.Attach(h1)
	set.Attach(h1)
	h1.AssertExpectations(t)
}

func TestHandlerSetEach(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	h1 := new(MockedHandler)
	h1.On("OnAttach").Return()
	h2 := new(MockedHandler)
	h2.On("OnAttach").Return()

	{
		hs := []Handler{}
		set.Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		set.Attach(h1)
		hs := []Handler{}
		set.Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, h1)
	}

	{
		set.Attach(h2)
		hs := []Handler{}
		set.Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, h1)
		assert.Contains(hs, h2)
	}
}

func TestHandlerSetEvery(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	h1 := new(MockedHandler)
	h1.On("OnAttach").Return()
	h2 := new(MockedHandler)
	h2.On("OnAttach").Return()

	assert.True(set.Every(func(h Handler) bool { return true }))
	assert.True(set.Every(func(h Handler) bool { return false }))

	set.Attach(h1)
	assert.True(set.Every(func(h Handler) bool { return true }))
	assert.False(set.Every(func(h Handler) bool { return false }))

	set.Attach(h2)
	assert.True(set.Every(func(h Handler) bool { return true }))
	assert.False(set.Every(func(h Handler) bool { return h.(*MockedHandler) == h1 }))
	assert.False(set.Every(func(h Handler) bool { return h.(*MockedHandler) == h2 }))
	assert.False(set.Every(func(h Handler) bool { return false }))
}

func TestHandlerSetSome(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	h1 := new(MockedHandler)
	h1.On("OnAttach").Return()
	h2 := new(MockedHandler)
	h2.On("OnAttach").Return()

	assert.False(set.Some(func(h Handler) bool { return true }))
	assert.False(set.Some(func(h Handler) bool { return false }))

	set.Attach(h1)
	assert.True(set.Some(func(h Handler) bool { return true }))
	assert.False(set.Some(func(h Handler) bool { return false }))

	set.Attach(h2)
	assert.True(set.Some(func(h Handler) bool { return true }))
	assert.True(set.Some(func(h Handler) bool { return h.(*MockedHandler) == h1 }))
	assert.True(set.Some(func(h Handler) bool { return h.(*MockedHandler) == h2 }))
	assert.False(set.Some(func(h Handler) bool { return false }))
}

func TestHandlerSetBoundEach(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedHandler)
	h.On("OnAttach").Return().Once()
	sh := new(MockedHandlerS)
	sh.On("OnAttach").Return().Once()
	oh := new(MockedHandlerO)
	oh.On("OnAttach").Return().Once()
	soh := new(MockedHandlerSO)
	soh.On("OnAttach").Return().Once()

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
		hs := []Handler{}
		set.BindSubject(s).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, sh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.BindSubject(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.BindSubject(new(Unit)).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.BindObject(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, oh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.BindObject(s).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.BindObject(new(Unit)).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.Bind(s).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, sh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.Bind(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 2)
		assert.Contains(hs, oh)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.Bind(new(Unit)).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.BindSubject(s).BindObject(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.BindObject(s).BindSubject(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Empty(hs)
	}

	{
		hs := []Handler{}
		set.BindSubject(s).Bind(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}

	{
		hs := []Handler{}
		set.Bind(s).BindObject(o).Each(func(h Handler) { hs = append(hs, h) })
		assert.Len(hs, 1)
		assert.Contains(hs, soh)
	}

	h.On("OnDetach").Return().Once()
	set.Bind(new(Unit)).Detach(h)
	h.AssertExpectations(t)
}

func TestHandlerSetBoundEvery(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedHandler)
	h.On("OnAttach").Return()
	sh := new(MockedHandlerS)
	sh.On("OnAttach").Return()
	sh.On("Subject").Return(s)
	oh := new(MockedHandlerO)
	oh.On("OnAttach").Return()
	oh.On("Object").Return(o)
	soh := new(MockedHandlerSO)
	soh.On("OnAttach").Return()
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)

	set.Attach(h)
	set.Attach(sh)
	set.Attach(oh)
	set.Attach(soh)

	assert.True(set.Bind(new(Unit)).Every(func(h Handler) bool { return true }))
	assert.True(set.Bind(new(Unit)).Every(func(h Handler) bool { return false }))

	assert.False(set.BindSubject(s).Every(func(h Handler) bool {
		assert.Equal(s, h.(Subject).Subject())
		if _, ok := h.(Object); ok {
			return h.(Object).Object() == o
		} else {
			return false
		}
	}))

	assert.False(set.BindObject(o).Every(func(h Handler) bool {
		assert.Equal(o, h.(Object).Object())
		if _, ok := h.(Subject); ok {
			return h.(Subject).Subject() == s
		} else {
			return false
		}
	}))
}

func TestHandlerSetBoundSome(t *testing.T) {
	assert := assert.New(t)
	set := MakeHandlerSet()
	s := new(Unit)
	o := new(Unit)
	h := new(MockedHandler)
	h.On("OnAttach").Return()
	sh := new(MockedHandlerS)
	sh.On("OnAttach").Return()
	sh.On("Subject").Return(s)
	oh := new(MockedHandlerO)
	oh.On("OnAttach").Return()
	oh.On("Object").Return(o)
	soh := new(MockedHandlerSO)
	soh.On("OnAttach").Return()
	soh.On("Subject").Return(s)
	soh.On("Object").Return(o)

	set.Attach(h)
	set.Attach(sh)
	set.Attach(oh)
	set.Attach(soh)

	assert.False(set.Bind(new(Unit)).Some(func(h Handler) bool { return true }))
	assert.False(set.Bind(new(Unit)).Some(func(h Handler) bool { return false }))

	assert.True(set.BindSubject(s).Some(func(h Handler) bool {
		assert.Equal(s, h.(Subject).Subject())
		if _, ok := h.(Object); ok {
			return h.(Object).Object() == o
		} else {
			return false
		}
	}))

	assert.True(set.BindObject(o).Some(func(h Handler) bool {
		assert.Equal(o, h.(Object).Object())
		if _, ok := h.(Subject); ok {
			return h.(Subject).Subject() == s
		} else {
			return false
		}
	}))
}
