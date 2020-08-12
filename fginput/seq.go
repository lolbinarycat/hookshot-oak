package fginput

func Seq(s string) []Input {
	ret := make([]Input,len(([]rune)(s)))
	var nxt Input
	for i, r := range ([]rune)(s) {
		switch r {
		case '←':
			nxt = Left
		case '↑':
			nxt = Up
		case '→':
			nxt = Right
		case '↓':
			nxt = Down
		case '↖':
			nxt = Up|Left
		case '↗':
			nxt = Up|Right
		case '↘':
			nxt = Down|Right
		case '↙':
			nxt = Down|Left
		case ' ':
			nxt = None
		default:
			nxt = Invalid
		}
		ret[i] = nxt
	}
	return ret
}
