package apcgen

type IntRange struct {
	Min int
	Max int
}

type keyValuePair[KT, VT any] struct {
	key   KT
	value VT
}
