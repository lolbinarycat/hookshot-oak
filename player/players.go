package player

// This file is bit of a hack

const MaxPlayers = 1

var players [1]*Player

func SetPlayer(i int,p *Player) {
	players[i] = p
}

func GetPlayer(i int) *Player {
	return players[i]
}
