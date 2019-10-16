package simdjson

import (
	"fmt"
	"sync"
	"testing"
)

func benchmarkFromFile(b *testing.B, filename string) {

	_, _, msg := loadCompressed(b, filename)

	b.SetBytes(int64(len(msg)))
	b.ReportAllocs()
	b.ResetTimer()

	pj := internalParsedJson{}
	pj.initialize(len(msg) * 2)

	for i := 0; i < b.N; i++ {
//		pj.structural_indexes = pj.structural_indexes[:0]
		pj.Tape = pj.Tape[:0]
		pj.Strings = pj.Strings[:0]
		find_structural_indices(msg, &pj)
		unified_machine(msg, &pj)
	}
}

func BenchmarkApache_builds(b *testing.B)  { benchmarkFromFile(b, "apache_builds") }
func BenchmarkCanada(b *testing.B)         { benchmarkFromFile(b, "canada") }
func BenchmarkCitm_catalog(b *testing.B)   { benchmarkFromFile(b, "citm_catalog") }
func BenchmarkGithub_events(b *testing.B)  { benchmarkFromFile(b, "github_events") }
func BenchmarkGsoc_2018(b *testing.B)      { benchmarkFromFile(b, "gsoc-2018") }
func BenchmarkInstruments(b *testing.B)    { benchmarkFromFile(b, "instruments") }
func BenchmarkMarine_ik(b *testing.B)      { benchmarkFromFile(b, "marine_ik") }
func BenchmarkMesh(b *testing.B)           { benchmarkFromFile(b, "mesh") }
func BenchmarkMesh_pretty(b *testing.B)    { benchmarkFromFile(b, "mesh.pretty") }
func BenchmarkNumbers(b *testing.B)        { benchmarkFromFile(b, "numbers") }
func BenchmarkRandom(b *testing.B)         { benchmarkFromFile(b, "random") }
func BenchmarkTwitter(b *testing.B)        { benchmarkFromFile(b, "twitter") }
func BenchmarkTwitterescaped(b *testing.B) { benchmarkFromFile(b, "twitterescaped") }
func BenchmarkUpdate_center(b *testing.B)  { benchmarkFromFile(b, "update-center") }

func testStage2Dev(t *testing.T, filename string) {

	_, _, msg := loadCompressed(t, filename)

	pj := internalParsedJson{}
	pj.initialize(len(msg)*2)

	pj.masks_chan = make(chan uint64, 1024)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		find_structural_indices(msg, &pj)
		wg.Done()
	}()

	go func() {

		var done bool
		i := uint32(0)     // index of the structural character (0,1,2,3...)
		idx := uint32(0)   // location of the structural character in the input (buf)
		c := byte(0)       // used to track the (structural) character we are looking at

		indexes := [64]uint32{} // make([]uint32, 0, 64)
		maskIndex := 0
		indexesLen := uint32(0)

		for {

			done, i, idx, c = UPDATE_CHAR_V3(msg, &pj, i, &indexes, &maskIndex, &indexesLen)
			fmt.Println(done, i, idx, string(c))
			if done {
				break
			}
		}

		wg.Done()
	}()

	wg.Wait()
}

func TestGscoDev(t *testing.T) { testStage2Dev(t, "gsoc-2018") }

func benchmarkStage2Dev(b *testing.B, filename string) {

	_, _, msg := loadCompressed(b, filename)

	b.SetBytes(int64(len(msg)))
	b.ReportAllocs()
	b.ResetTimer()

	pj := internalParsedJson{}
	pj.initialize(len(msg)*2)

	// buf := make([]byte, len(pj.masks)*8)
	// for i, m := range pj.masks {
	// 	binary.LittleEndian.PutUint64(buf[i*8:], m)
	// }
	// ioutil.WriteFile(filename+".masks", []byte(buf), 0666)

	for i := 0; i < b.N; i++ {

		var wg sync.WaitGroup
		wg.Add(2)

		pj.masks_chan = make(chan uint64, 1024*1024)

		go func() {
			find_structural_indices(msg, &pj)
			wg.Done()
		}()

		go func() {
			pj.Tape = pj.Tape[:0]
			pj.Strings = pj.Strings[:0]
			unified_machine(msg, &pj)
			wg.Done()
		}()

		wg.Wait()
	}
}

// BenchmarkGsocDev-8          2000            939894 ns/op        3540.64 MB/s       43270 B/op          0 allocs/op
// BenchmarkGsocDev-8          2000            930941 ns/op        3574.69 MB/s       43270 B/op          0 allocs/op
// BenchmarkGsocDev-8          2000            961418 ns/op        3461.38 MB/s       43270 B/op          0 allocs/op

func BenchmarkGsocDev(b *testing.B) { benchmarkStage2Dev(b, "gsoc-2018") }
