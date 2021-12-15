package pool

import (
	"container/list"
	"fmt"
	"lms-db/strategy"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type A struct {
	list  *list.List
	mutex sync.RWMutex
}

type BB struct {
	A     string
	Mutex sync.Mutex
}

func TestCD(t *testing.T) {
	a := strategy.NewLruStrategy(10)

	b1 := &BB{
		A: "a",
	}

	b2 := &BB{
		A: "b",
	}

	b3 := &BB{
		A: "c",
	}

	a.Set("a", b1, nil)
	a.Set("b", b2, nil)
	a.Set("c", b3, nil)

	go func() {
		if c, ok := a.Get("a"); ok {
			if ai, ok := c.Value.(*BB); ok {
				ai.Mutex.Lock()
				time.Sleep(time.Second * 10)
				ai.Mutex.Unlock()
			}
		}
	}()
	time.Sleep(time.Second)

	if c, ok := a.Get("a"); ok {
		if ai, ok := c.Value.(*BB); ok {
			fmt.Println(ai)
			fmt.Println(1231231313131)
		}
	} else {
		fmt.Println(ok)
	}

}

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
	fmt.Println(filepath.Join("/123", "a"))

	var a int = 10
	func() {
		fmt.Println(a)
		a = 20
	}()
	fmt.Println(a)
}

func TestAE(t *testing.T) {
	var number float32 = 0.085
	fmt.Printf("Starting number:%f\n", number)
	bits := math.Float32bits(number)
	binary := fmt.Sprintf("%.32b", bits)
	bias := 127
	_ = bits & (1 << 31)
	expointentRaw := int(bits >> 23)
	expointent := expointentRaw - bias
	var mantissa float64
	a := binary[9:32]
	c := []byte(a)
	fmt.Println(a, c)
	for index, bit := range binary[9:32] {
		if bit == 49 {
			position := index + 1
			bitValue := math.Pow(2, float64(position))
			fractional := 1 / bitValue
			mantissa = mantissa + fractional
		}
	}
	value := (1 + mantissa) * math.Pow(2, float64(expointent))
	fmt.Println(value)

}

func GA() chan int {
	ch := make(chan int, 10)
	fmt.Println("一直在调用")
	go func() {
		for {
			ch <- rand.Int()
		}
	}()
	return ch
}

func TestCDEDFF(t *testing.T) {
	GAch := GA()
	ch := make(chan int, 20)
	go func() {
		for {
			select {
			case ch <- <-GAch:
			case <-time.Tick(30 * time.Second):
			case <-time.After(30 * time.Second):

			}
		}
	}()
	//
	for {
		select {
		case r := <-ch:
			fmt.Println("读出数据来了", r)
		}
	}
}

func GB() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <- rand.Int()
			fmt.Println("1")
		}
	}()
	fmt.Println("b")
	return ch
}

func TestCDEDF(t *testing.T) {
	ch := make(chan int, 20)
	go func() {
		select {
		case ch <- <-GA():
			fmt.Println("A=======", ch)
		case ch <- <-GB():
			fmt.Println("B=======", ch)

		}
	}()
	time.Sleep(time.Second * 12)
}

func TestCDA(t *testing.T) {
	var wg sync.WaitGroup
	foo := make(chan int)
	bar := make(chan int)
	closing := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case foo <- <-bar:
		case <-closing:
			println("closing")
		}
	}()
	// bar <- 123
	close(closing)
	wg.Wait()
}

func talk(msg string, sleep int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}()
	return ch
}

func fanIn(input1, input2 <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case ch <- <-input1:
			case ch <- <-input2:
			}
		}
	}()
	return ch
}

func TestCa(t *testing.T) {
	ch := fanIn(talk("A", 10), talk("B", 1000))
	for i := 0; i < 10; i++ {
		fmt.Printf("%q\n", <-ch)
	}
}

func TestTAADSADA(t *testing.T) {
	ch := make(chan string)
	//go func() {
	//	for {
	//		//ch <- "y"
	//	}
	//}()
	go func() {
		if err := http.ListenAndServe("127.0.0.1:9999", nil); err != nil {
			fmt.Println("http server fail")
		}
	}()
	//tk := time.Tick(time.Millisecond)
	for {
		select {
		case <-ch:
		case <-time.Tick(time.Hour):
		}
	}

}
