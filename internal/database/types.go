package database

type List struct {
	items []string
}

type Hash struct {
	fields map[string]string
}

type DataType int

const (
	String DataType = iota
	ListType
	HashType
)
