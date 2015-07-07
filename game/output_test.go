package game

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
		expected := []byte(`{"Type":"Level","Data":{}}`)
		actual, err := EncodeOutputFrame(OutputLevel{})
		if assert.Nil(err) {
			assert.Equal(expected, actual)
		}
	}

	{
		expected := []byte(`{"Type":"Player","Data":{"UnitID":100,"Q":{"Name":"Q","Description":"Q Description","TargetType":3,"HealthCost":0,"ManaCost":10,"ActivationDuration":0,"CooldownDuration":40,"DisableTypes":["Stun"]},"W":{"Name":"W","Description":"W Description","TargetType":3,"HealthCost":0,"ManaCost":20,"ActivationDuration":0,"CooldownDuration":160,"DisableTypes":["Stun"]},"E":{"Name":"E","Description":"E Description","TargetType":3,"HealthCost":0,"ManaCost":40,"ActivationDuration":0,"CooldownDuration":400,"DisableTypes":["Stun"]},"R":{"Name":"R","Description":"R Description","TargetType":3,"HealthCost":0,"ManaCost":80,"ActivationDuration":0,"CooldownDuration":1200,"DisableTypes":["Stun"]}}}`)
		actual, err := EncodeOutputFrame(OutputPlayer{
			UnitID: 100,
			Q: OutputPlayerAbility{
				Name:               "Q",
				Description:        "Q Description",
				TargetType:         TargetTypeEnemy,
				HealthCost:         0,
				ManaCost:           10,
				ActivationDuration: 0,
				CooldownDuration:   2 * Second,
				DisableTypes:       []string{"Stun"},
			},
			W: OutputPlayerAbility{
				Name:               "W",
				Description:        "W Description",
				TargetType:         TargetTypeEnemy,
				HealthCost:         0,
				ManaCost:           20,
				ActivationDuration: 0,
				CooldownDuration:   8 * Second,
				DisableTypes:       []string{"Stun"},
			},
			E: OutputPlayerAbility{
				Name:               "E",
				Description:        "E Description",
				TargetType:         TargetTypeEnemy,
				HealthCost:         0,
				ManaCost:           40,
				ActivationDuration: 0,
				CooldownDuration:   20 * Second,
				DisableTypes:       []string{"Stun"},
			},
			R: OutputPlayerAbility{
				Name:               "R",
				Description:        "R Description",
				TargetType:         TargetTypeEnemy,
				HealthCost:         0,
				ManaCost:           80,
				ActivationDuration: 0,
				CooldownDuration:   60 * Second,
				DisableTypes:       []string{"Stun"},
			},
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
