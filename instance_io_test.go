package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceInput(t *testing.T) {
	assert := assert.New(t)
	w := MakeInstanceInput(1000)
	assert.Implements((*InstanceInputWriter)(nil), w)

	w.Write(1, InputChat{
		Message: "hi",
	})
	assert.NotPanics(func() {
		input := <-w
		assert.Equal(ClientID(1), input.ClientID)
		assert.IsType(InputChat{}, input.Input)
	})
}

func TestInstanceOutput(t *testing.T) {
	assert := assert.New(t)
	w := MakeInstanceOutput(1000)
	assert.Implements((*InstanceOutputWriter)(nil), w)

	w.Write(OutputChat{
		ClientName: "user",
		Message:    "hi all",
	})
	assert.NotPanics(func() {
		output := <-w
		assert.Equal(ClientID(0), output.ClientID)
		assert.IsType(OutputChat{}, output.Output)
	})

	w.BindClientID(1).Write(OutputMessage{
		Message: "welcome",
	})
	assert.NotPanics(func() {
		output := <-w
		assert.Equal(ClientID(1), output.ClientID)
		assert.IsType(OutputMessage{}, output.Output)
	})

	w.BindClientID(1).BindClientID(2).Write(OutputMessage{
		Message: "welcome",
	})
	assert.NotPanics(func() {
		output := <-w
		assert.Equal(ClientID(2), output.ClientID)
		assert.IsType(OutputMessage{}, output.Output)
	})
}
