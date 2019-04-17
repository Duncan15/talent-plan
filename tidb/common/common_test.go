package main

import "testing"

//import (
//	"TiDB/alg/sort"
//	"TiDB/gene"
//	"TiDB/record"
//	"log"
//)
//
//func main()  {
//	log.Print("start benchmark")
//	for i := 0; i < 3; i++ {
//		arr1 := gene.GenerateRandomInt(100000, 10)
//		log.Printf("slice size is %d", len(arr1))
//		arr2 := make([]int, len(arr1))
//		arr3 := make([]int, len(arr1))
//		arr4 := make([]int, len(arr1))
//		copy(arr2, arr1)
//		copy(arr3, arr1)
//		copy(arr4, arr1)
//		//record.RecordTime("customMergeSort", func() {
//		//	sort.MergeSort(arr2, 0, len(arr2) - 1)
//		//})
//		//record.RecordTime("customQuickSort", func() {
//		//	sort.QuickSort(arr3, 0, len(arr3) - 1)
//		//})
//		//record.RecordTime("standardQuickSort", func() {
//		//	sort2.Ints(arr1)
//		//})
//
//		record.RecordTime("mergeSortByMultiGoroutine", func() {
//			sort.MergeSortByMultiGoroutine(arr4, 0, len(arr4) - 1)
//		})
//	}
//
//}
func TestCommon(t *testing.T)  {

}