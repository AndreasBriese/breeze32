// breeze32_test.go

package breeze32

import (
	"fmt"
	"math/rand"
	// "strings"
	"os"
	// "sync"
	"testing"
	// "time"
)

const CHARS = "!abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890#"

const SAMPLELEN = 1000000

var prng = rand.New(rand.NewSource(int64(12345))) //time.Now().Nanosecond())))
var randb = new(Breeze32).Init(12345)

func exportCsv(fleName string) chan int {
	fle, err := os.Create(fleName)
	if err != nil {
		panic(err)
	}
	ch := make(chan int)
	go func(ch chan int, fle *os.File) {
	loop:
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					break loop
				}
				fle.WriteString(fmt.Sprintf("%v,\n", v))
			}
		}
		fle.Close()
	}(ch, fle)
	return ch
}

func Benchmark_Init_and_Seed(b *testing.B) {
	for r := 0; r < b.N; r++ {
		rand := new(Breeze32)
		rd := byte(0)
		rand.Init(123).Byte(&rd)
		rd = rd
	}
}

// func Benchmark_breezeByte(b *testing.B) {
// 	// randb := new(Breeze32)
// 	// randb.Init(123)
// 	rd := byte(0)
// 	b.ResetTimer()
// 	for r := 0; r < b.N; r++ {
// 		randb.Byte(&rd)
// 		rd = rd
// 	}
// }

func Benchmark_breeze_Uint64(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	rd := uint64(0)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd = randb.Uint64(&rd)
	}
}

func Benchmark_breeze_Uint32(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	rd := uint32(0)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd = randb.Uint32(&rd)
	}
}

// func Benchmark_breezeB16(b *testing.B) {
// 	randb := new(Breeze32)
// 	rd := uint16(0)
// 	randb.Init(123)
// 	b.ResetTimer()
// 	for r := 0; r < b.N; r++ {
// 		rd = randb.Uint16(&rd)
// 	}
// }

func Benchmark_breeze_Uint8(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	rd := uint8(0)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd = randb.Uint8(&rd)
	}
}

func Benchmark_breeze_ByteMP(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	rd := uint8(0)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd = randb.ByteMP(&rd)
	}
}

func Benchmark_breeze_RandIntN_8byte(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := randb.RandIntN(1 << 7)
		rd = rd
	}
}

func Benchmark_breeze_RandIntN_16byte(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := randb.RandIntN(1 << 12)
		rd = rd
	}
}

func Benchmark_breeze_RandIntN_32byte(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := randb.RandIntN(1 << 32)
		rd = rd
	}
}

func Benchmark_breeze_RandIntN_64byte(b *testing.B) {
	// randb := new(Breeze32)
	// randb.Init(123)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := randb.RandIntN(1 << 62)
		rd = rd
	}
}

// func Benchmark_mathRand256(b *testing.B) {
// 	for r := 0; r < b.N; r++ {
// 		rd := prng.Intn(256)
// 		rd = rd
// 	}
// }

func Benchmark_mathRand_Intn_8byte(b *testing.B) {
	prng = rand.New(rand.NewSource(int64(123)))
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := prng.Intn(1 << 7)
		rd = rd
	}
}

func Benchmark_mathRand_Intn_16byte(b *testing.B) {
	prng = rand.New(rand.NewSource(int64(123)))
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := prng.Intn(1 << 12)
		rd = rd
	}
}

func Benchmark_mathRand_Intn_32byte(b *testing.B) {
	prng = rand.New(rand.NewSource(int64(123)))
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := prng.Intn(1 << 32)
		rd = rd
	}
}

func Benchmark_mathRand_Intn_64byte(b *testing.B) {
	prng = rand.New(rand.NewSource(int64(123)))
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		rd := prng.Intn(1<<63 - 1)
		rd = rd
	}
}

/*
func Example_CsVbreezeRandIntN65535() {
	randb := new(Breeze32)
	rd := byte(0)
	randb.Init(123)
	ch := exportCsv("breezeRand.csv")
	sum := 0
	for i := 0; i < 1000000; i++ {
		ch <- int(randb.Byte(&rd))
		sum += int(rd)
	}
	close(ch)
	fmt.Println(sum)
	// Output: lhöjkhjkl
}
*/
/*
func Example_CsVPrngIntN65535() {
	rand := new(Breeze32)
	rand.Init(123)
	sum := 0
	ch := exportCsv("mathRand.csv")
	for i := 0; i < 1000000; i++ {
		rd := prng.Intn(1 << 12)
		sum += rd
		ch <- rd
	}
	close(ch)
	fmt.Println(sum)
	// Output: ljlkjlöj
}

*/
// func Example_CsVbreezeRandIntN65535() {
// 	randb := new(Breeze32)
// 	// randb.Init(uint64(time.Now().UnixNano()))
// 	randb.Init(123456789)
// 	sum := 0
// 	ch := exportCsv("breezeRand.csv")
// 	for i := 0; i < 1000000; i++ {
// 		rd := randb.RandIntN(166)
// 		sum += rd
// 		ch <- rd
// 	}
// 	close(ch)
// 	fmt.Println(sum)
// 	// Output: 82462458
// }

/*

func Example_CsVPrngIntN65535() {
	rand := new(Breeze32)
	rand.Init(123456789)
	sum := 0
	ch := exportCsv("mathRand.csv")
	for i := 0; i < 1000000; i++ {
		rd := prng.Intn(1 << 35) // rand.RandIntN(65535)
		sum += rd
		ch <- rd
	}
	close(ch)
	fmt.Println(sum)
	// Output: 4502622
}
*/

// func ExampleMP() {
// 	randb := new(Breeze32)
// 	randb.Init(123)
// 	wg := new(sync.WaitGroup)

// 	mpF := func(randb *Breeze32, wg *sync.WaitGroup) {
// 		l := make([]uint8, 1000000)
// 		r := uint8(0)
// 		for i := range l {
// 			l[i] = randb.ByteMP(&r)
// 		}
// 		// fmt.Println(l)
// 		wg.Done()
// 	}

// 	wg.Add(12)
// 	st := time.Now()
// 	// 4 Goroutines
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	go mpF(randb, wg)
// 	wg.Wait()
// 	fmt.Println(time.Since(st).Nanoseconds()/12000000, "ns/op 12 parallel (total 12e6)")
// 	////  Output: ff
// }
