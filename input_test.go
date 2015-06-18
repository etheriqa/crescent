package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeInputFrame(t *testing.T) {
	assert := assert.New(t)

	{
		json := []byte(``)
		_, err := DecodeInputFrame(json)
		assert.NotNil(err)
	}

	{
		json := []byte(`{"Type":"UnknownType","Data":{}}`)
		_, err := DecodeInputFrame(json)
		assert.NotNil(err)
	}

	{
		json := []byte(`{"Type":"Interrupt"}`)
		_, err := DecodeInputFrame(json)
		assert.NotNil(err)
	}

	{
		json := []byte(`{"Type":"Interrupt","Data":null}`)
		_, err := DecodeInputFrame(json)
		assert.NotNil(err)
	}

	{
		json := []byte(`{"Type":"Profile","Data":{"UserName":"etheriqa"}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputProfile{}, input) {
			assert.Equal(UserName("etheriqa"), input.(InputProfile).UserName)
		}
	}

	{
		json := []byte(`{"Type":"Chat","Data":{"Message":"hi all"}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputChat{}, input) {
			assert.Equal("hi all", input.(InputChat).Message)
		}
	}

	{
		json := []byte(`{"Type":"Stage","Data":{"StageID":1}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputStage{}, input) {
			assert.Equal(StageID(1), input.(InputStage).StageID)
		}
	}

	{
		json := []byte(`{"Type":"Join","Data":{"ClassName":"Healer"}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputJoin{}, input) {
			assert.Equal(ClassName("Healer"), input.(InputJoin).ClassName)
		}
	}

	{
		json := []byte(`{"Type":"Leave","Data":{}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputLeave{}, input) {
		}
	}

	{
		json := []byte(`{"Type":"Ability","Data":{"AbilityName":"Q","ObjectUnitID":null}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputAbility{}, input) {
			assert.Equal("Q", input.(InputAbility).AbilityName)
			assert.Nil(input.(InputAbility).ObjectUnitID)
		}
	}

	{
		json := []byte(`{"Type":"Ability","Data":{"AbilityName":"W","ObjectUnitID":1}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputAbility{}, input) {
			assert.Equal("W", input.(InputAbility).AbilityName)
			if assert.NotNil(input.(InputAbility).ObjectUnitID) {
				assert.Equal(UnitID(1), *input.(InputAbility).ObjectUnitID)
			}
		}
	}

	{
		json := []byte(`{"Type":"Interrupt","Data":{}}`)
		input, err := DecodeInputFrame(json)
		if assert.Nil(err) && assert.IsType(InputInterrupt{}, input) {
		}
	}
}
