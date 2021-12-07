package pool

import (
	"fmt"
	"testing"
)

func Test_A(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(a[0:2])
	fmt.Println(hash([]byte("abc")))
	fmt.Println(hash([]byte("abcd")))

}

func hash(key []byte) uint {
	var nr, nr2 uint = 1, 4
	for i := len(key); i > 0; i-- {
		nr ^= ((nr & 63) + nr2) * (uint(key[i-1]) + (nr << 8))
		nr2 += 3
	}
	return nr
}

type StructA struct {
	A string
}

func (a *StructA) PrintA() {
	a.A = "121231"
}

func TestC(t *testing.T) {
	a := StructA{A: "A"}
	b := StructA{A: "C"}
	a.PrintA()
	fmt.Println(a)
	fmt.Println(b)
}

func TestD(t *testing.T) {
	a := "\n"
	fmt.Println([]byte(a))

	b := []byte{1, 2}
	fmt.Println(b[0:2])
	count := 5
	for i := count; i > 0; i-- {
		fmt.Println(i)
	}
}

func TestE(t *testing.T) {
	var a int = 10
	func() {
		fmt.Println(a)
		a = 20
	}()
	fmt.Println(a)
}
