package main

import "pingcap/talentplan/tidb/common/alg/sort"

//import "pingcap/talentplan/tidb/common/alg/sort"

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	sort.MergeSortByMultiGoroutine(src, 0, len(src) - 1)
}
