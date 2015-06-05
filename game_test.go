package main

func mockGame() *game {
	return newGame(make(chan message, 1000), make(chan message, 1000))
}
