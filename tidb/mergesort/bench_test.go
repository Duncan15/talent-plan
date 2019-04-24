package main

import (
	"sort"
	"testing"
	cussort "pingcap/talentplan/tidb/common/alg/sort"
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
		cussort.Merge(original, 0 ,len(original) / 2, len(original) - 1, src)
		b.StopTimer()
	}
}

func BenchmarkMergeSort(b *testing.B) {
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
