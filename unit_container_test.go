package crescent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitMap(t *testing.T) {
	assert := assert.New(t)

	us := NewUnitMap()
	assert.Implements((*UnitContainer)(nil), us)
	assert.Implements((*UnitQueryable)(nil), us)

	assert.NotNil(us.Leave(1))

	players := make([]*Unit, 5)
	for i := 0; i < GroupCapacity; i++ {
		u, err := us.Join(UnitGroupPlayer, "Player", new(Class))
		assert.NotNil(u)
		assert.Nil(err)
		players[i] = u
	}
	{
		u, err := us.Join(UnitGroupPlayer, "Failure", new(Class))
		assert.Nil(u)
		assert.NotNil(err)
	}

	for _, player := range players {
		assert.Equal(player, us.Find(player.ID()))
	}

	assert.Nil(us.Leave(players[0].ID()))
	assert.Nil(us.Find(players[0].ID()))
	assert.NotNil(us.Leave(players[0].ID()))

	us.Clear()
	for _, player := range players {
		assert.Nil(us.Find(player.ID()))
	}
}

func TestUnitMapEach(t *testing.T) {
	assert := assert.New(t)
	query := NewUnitMap()

	{
		units := []*Unit{}
		query.Each(func(u *Unit) { units = append(units, u) })
		assert.Empty(units)
	}

	p1, _ := query.Join(UnitGroupPlayer, "Player1", new(Class))

	{
		units := []*Unit{}
		query.Each(func(u *Unit) { units = append(units, u) })
		assert.Len(units, 1)
		assert.Contains(units, p1)
	}

	{
		friends := []*Unit{}
		query.EachFriend(p1, func(u *Unit) { friends = append(friends, u) })
		assert.Len(friends, 1)
		assert.Contains(friends, p1)
	}

	{
		enemies := []*Unit{}
		query.EachEnemy(p1, func(u *Unit) { enemies = append(enemies, u) })
		assert.Empty(enemies)
	}

	a1, _ := query.Join(UnitGroupAI, "AI1", new(Class))

	{
		units := []*Unit{}
		query.Each(func(u *Unit) { units = append(units, u) })
		assert.Len(units, 2)
		assert.Contains(units, p1)
		assert.Contains(units, a1)
	}

	{
		friends := []*Unit{}
		query.EachFriend(p1, func(u *Unit) { friends = append(friends, u) })
		assert.Len(friends, 1)
		assert.Contains(friends, p1)
	}

	{
		enemies := []*Unit{}
		query.EachEnemy(p1, func(u *Unit) { enemies = append(enemies, u) })
		assert.Len(enemies, 1)
		assert.Contains(enemies, a1)
	}

	{
		friends := []*Unit{}
		query.EachFriend(a1, func(u *Unit) { friends = append(friends, u) })
		assert.Len(friends, 1)
		assert.Contains(friends, a1)
	}

	{
		enemies := []*Unit{}
		query.EachEnemy(a1, func(u *Unit) { enemies = append(enemies, u) })
		assert.Len(enemies, 1)
		assert.Contains(enemies, p1)
	}

	p2, _ := query.Join(UnitGroupPlayer, "Player2", new(Class))
	a2, _ := query.Join(UnitGroupAI, "AI2", new(Class))

	{
		units := []*Unit{}
		query.Each(func(u *Unit) { units = append(units, u) })
		assert.Len(units, 4)
		assert.Contains(units, p1)
		assert.Contains(units, p2)
		assert.Contains(units, a1)
		assert.Contains(units, a2)
	}

	{
		friends := []*Unit{}
		query.EachFriend(p1, func(u *Unit) { friends = append(friends, u) })
		assert.Len(friends, 2)
		assert.Contains(friends, p1)
		assert.Contains(friends, p2)
	}

	{
		enemies := []*Unit{}
		query.EachEnemy(p1, func(u *Unit) { enemies = append(enemies, u) })
		assert.Len(enemies, 2)
		assert.Contains(enemies, a1)
		assert.Contains(enemies, a2)
	}
}
