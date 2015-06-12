package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeUnitPair(t *testing.T) {
	assert := assert.New(t)
	subject := new(Unit)
	object := new(Unit)

	assert.Panics(func() { MakeUnitPair(nil, nil) })
	assert.Panics(func() { MakeUnitPair(nil, object) })
	assert.Panics(func() { MakeUnitPair(subject, nil) })
	up := MakeUnitPair(subject, object)
	assert.Equal(subject, up.Subject())
	assert.Equal(object, up.Object())
}

func TestMakeSubject(t *testing.T) {
	assert := assert.New(t)
	subject := new(Unit)

	assert.Panics(func() { MakeSubject(nil) })
	up := MakeSubject(subject)
	assert.Equal(subject, up.Subject())
	assert.Nil(up.Object())
}

func TestMakeObject(t *testing.T) {
	assert := assert.New(t)
	object := new(Unit)

	assert.Panics(func() { MakeObject(nil) })
	up := MakeObject(object)
	assert.Nil(up.Subject())
	assert.Equal(object, up.Object())
}
