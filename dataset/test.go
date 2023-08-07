package dataset

import "C"

type GenHeader[T any] struct {
	List []T
}

func (g *GenHeader[T]) SetGI(i int, x T) {
	g.List[i] = x
}

func (g *GenHeader[T]) GetGI(i int) T {
	return g.List[i]
}
func (v *Label) SetLines(lines int) {
	C.gtk_label_set_lines(v.native(), C.gint(lines))
}
func (v *Grid) RemoveColumn(position int) {
	C.gtk_grid_remove_column(v.native(), C.gint(position))
}
