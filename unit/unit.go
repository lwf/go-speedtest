package unit

const (
	Bi   = iota
	Kibi = 1 << (10 * iota)
	Mebi
	Gibi
)

type Labels map[int]string

var Bits = Labels{
	Bi:   "bit",
	Kibi: "kbit",
	Mebi: "mbit",
	Gibi: "gbit",
}

var Bytes = Labels{
	Bi:   "b",
	Kibi: "kb",
	Mebi: "mb",
	Gibi: "gb",
}

func Unit(v float64, labels map[int]string) (float64, string) {
	switch {
	case v > Gibi:
		return v / Gibi, labels[Gibi]
	case v > Mebi:
		return v / Mebi, labels[Mebi]
	case v > Kibi:
		return v / Kibi, labels[Kibi]
	default:
		return v, labels[Bi]
	}
}
