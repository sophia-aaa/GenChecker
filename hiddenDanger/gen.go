package hiddenDanger

type GenHeader[T any] struct {
	List []T
}

func (g *GenHeader[T]) Gens() []T {
	return g.List
}

//func main() {}

/*


func (g *GenHeader[T]) SetGI(i int, x T) {
	g.List[i] = x
}

func (g *GenHeader[T]) GetGI(i int) T {
	return g.List[i]
}

func main() {

	aG := GenHeader[int]{List: []int{10, 11}}
	bG := GenHeader[int]{List: []int{0, 1, 2, 3, 4}}
	fmt.Println("aG is", aG)
	fmt.Println("bG is", bG)
	fmt.Println("aG.Gen() is", aG.Gens())
	fmt.Println("bG.Gen() is", bG.Gens())
	aG.SetGI(0, 0)
	fmt.Println("aG.Gen() is", aG.Gens())
	fmt.Println(aG.GetGI(1))
	fmt.Println()

	cG := GenHeader[bool]{List: []bool{false, false}}
	dG := GenHeader[string]{List: []string{"abs ", "cd", "ef", "ge", "wg2"}}
	fmt.Println("cG.Gen() is", cG.Gens())
	fmt.Println("dG.Gen() is", dG.Gens())
	cG.SetGI(0, true)
	dG.SetGI(1, "test")
	fmt.Println("cG.Gen() is", cG.Gens())
	fmt.Println("dG.Gen() is", dG.Gens())
	fmt.Println(cG.GetGI(0))
	fmt.Println()

	var var1 complex64 = complex(3, -5)
	var var2 complex64 = complex(2, 7.5)
	var var3 complex64 = complex(2, 4)

	eG := GenHeader[complex64]{List: []complex64{var1, var2}}
	fG := GenHeader[float32]{List: []float32{1.234, 2.54, 3.4, 5.6}}
	fmt.Println("eG.Gen() is", eG.Gens())
	fmt.Println("fG.Gen() is", fG.Gens())
	eG.SetGI(0, var3)
	fG.SetGI(1, 6.0)
	fmt.Println("eG.Gen() is", eG.Gens())
	fmt.Println("fG.Gen() is", fG.Gens())
	fmt.Println(eG.GetGI(0))
	fmt.Println()

	hG := GenHeader[*complex64]{List: []*complex64{&var1, &var2}}
	fmt.Println("hG.Gen() is", hG.Gens())
	hG.SetGI(0, &var3)
	fmt.Println("hG.Gen() is", hG.Gens())
	fmt.Println(hG.GetGI(0))
	fmt.Println()

	pG := GenHeader[uintptr]{List: []uintptr{uintptr(1), uintptr(2)}}
	fmt.Println("pG.Gen() is", pG.Gens())
	pG.SetGI(0, uintptr(4))
	fmt.Println("pG.Gen() is", pG.Gens())
	fmt.Println(pG.GetGI(0))
	fmt.Println()
}
*/
