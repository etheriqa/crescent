package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeOutputFrame(t *testing.T) {
	assert := assert.New(t)

	{
		_, err := EncodeOutputFrame(nil)
		assert.NotNil(err)
	}

	{
		expected := []byte(`{"Type":"Sync","Data":{"InstanceTime":2000}}`)
		actual, err := EncodeOutputFrame(OutputSync{
			InstanceTime: 2000,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"Message","Data":{"Message":"welcome"}}`)
		actual, err := EncodeOutputFrame(OutputMessage{
			Message: "welcome",
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"Chat","Data":{"UserName":"user","Message":"hi all"}}`)
		actual, err := EncodeOutputFrame(OutputChat{
			UserName: "user",
			Message:  "hi all",
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitJoin","Data":{"UnitID":100,"UnitGroup":0,"UnitName":"user","ClassName":"Healer","Health":500,"HealthMax":1000,"Mana":200,"ManaMax":400}}`)
		actual, err := EncodeOutputFrame(OutputUnitJoin{
			UnitID:    100,
			UnitGroup: 0,
			UnitName:  "user",
			ClassName: "Healer",
			Health:    500,
			HealthMax: 1000,
			Mana:      200,
			ManaMax:   400,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitLeave","Data":{"UnitID":100}}`)
		actual, err := EncodeOutputFrame(OutputUnitLeave{
			UnitID: 100,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitAttach","Data":{"UnitID":100,"AttachmentName":"Stun","Stack":0,"ExpirationTime":100}}`)
		actual, err := EncodeOutputFrame(OutputUnitAttach{
			UnitID:         100,
			AttachmentName: "Stun",
			Stack:          0,
			ExpirationTime: 100,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitDetach","Data":{"UnitID":100,"AttachmentName":"Stun"}}`)
		actual, err := EncodeOutputFrame(OutputUnitDetach{
			UnitID:         100,
			AttachmentName: "Stun",
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitActivating","Data":{"UnitID":100,"AbilityName":"Q","StartTime":2000,"EndTime":2020}}`)
		actual, err := EncodeOutputFrame(OutputUnitActivating{
			UnitID:      100,
			AbilityName: "Q",
			StartTime:   2000,
			EndTime:     2020,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitActivated","Data":{"UnitID":100,"AbilityName":"Q","OK":true}}`)
		actual, err := EncodeOutputFrame(OutputUnitActivated{
			UnitID:      100,
			AbilityName: "Q",
			OK:          true,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitCooldown","Data":{"UnitID":100,"AbilityName":"Q","ExpirationTime":3000,"Active":true}}`)
		actual, err := EncodeOutputFrame(OutputUnitCooldown{
			UnitID:         100,
			AbilityName:    "Q",
			ExpirationTime: 3000,
			Active:         true,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"UnitResource","Data":{"UnitID":100,"Health":1000,"Mana":400}}`)
		actual, err := EncodeOutputFrame(OutputUnitResource{
			UnitID: 100,
			Health: 1000,
			Mana:   400,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"Damage","Data":{"SubjectUnitID":100,"ObjectUnitID":101,"Damage":20,"IsCritical":false}}`)
		actual, err := EncodeOutputFrame(OutputDamage{
			SubjectUnitID: 100,
			ObjectUnitID:  101,
			Damage:        20,
			IsCritical:    false,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"Healing","Data":{"SubjectUnitID":100,"ObjectUnitID":101,"Healing":20,"IsCritical":false}}`)
		actual, err := EncodeOutputFrame(OutputHealing{
			SubjectUnitID: 100,
			ObjectUnitID:  101,
			Healing:       20,
			IsCritical:    false,
		})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}
}
