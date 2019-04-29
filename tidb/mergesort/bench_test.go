package main

import (
	cussort "pingcap/talentplan/tidb/common/alg/sort"
	"sort"
	"sync"
	"testing"
)

func compareAndSet(src []int64, original []int64)  {
	i, j := 0, len(original) / 2
	for index := 0; index < len(original); index++ {
		if index % 2 == 0 {
			src[i] = original[index]
			i++
		} else {
			src[j] = original[index]
			j++
		}
	}
}
func sequentialSet(src []int64, original []int64)  {
	for index := 0; index < len(original); index++ {
		src[index] = original[index]
	}
}
func BenchmarkVer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		numElements := 16 << 20
		src := make([]int64, numElements)
		original := make([]int64, numElements)
		prepare(original)
		b.StartTimer()
		l := len(original) / 10
		wg := &sync.WaitGroup{}
		for j, cur, end := 0, 0, 0; j < 10; j++ {
			end = cur + l
			if end > len(original) - 1 {
				end = len(original) - 1
			}
			wg.Add(1)
			go func(a int, b int) {
				cussort.Merge(original, a, (a + b) / 2, b, src)
				wg.Done()
			}(cur, end)
			cur = end + 1
		}
		wg.Wait()
		b.StopTimer()
	}
}

func BenchmarkMergeSort(b *testing.B) {
	//runtime.GOMAXPROCS(1)
	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		MergeSort(src)
	}
}

func BenchmarkNormalSort(b *testing.B) {
	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
	}
}
