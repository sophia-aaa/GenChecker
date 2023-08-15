package hiddenDanger

import "C"

type GenHeader[T any] struct {
	List []T
}

func (g *GenHeader[T]) Object() []T {
	return g.List
}

func (g *GenHeader[T]) SetGI(i int, x T) {
	g.List[i] = x
}

/*
func (g GenHeader[T]) GetGI(i int) T {
	return g.List[i]
}
*/
