package sort

import (
	"sync"
)

//MergeSort the export method for user
func MergeSort(arr []int64, start int, end int) {
	tmp := make([]int64, len(arr))
	mergeSort(arr, start, end, tmp)
}

func MergeSortByMultiGoroutine(arr []int64, start int, end int)  {
	tmp := make([]int64, len(arr))
	//mergeLen := (end - start + 1) / runtime.NumCPU()
	mergeSortByMultiGoroutine(arr, start, end, 1000, tmp, nil)
}
func mergeSortByMultiGoroutine(arr []int64, start int, end int, mergeLen int, tmp []int64, wg *sync.WaitGroup)  {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if end <= start {
		return
	}
	if end - start + 1 <= mergeLen {
		QuickSort(arr, start, end)
	} else {
		mid := start + (end - start) / 2
		twg := &sync.WaitGroup{}
		twg.Add(2)
		go mergeSortByMultiGoroutine(arr, start, mid, mergeLen, tmp, twg)
		go mergeSortByMultiGoroutine(arr, mid + 1, end, mergeLen, tmp, twg)
		twg.Wait()
		Merge(arr, start, mid, end, tmp)
	}
}

func mergeSort(arr []int64, start int, end int, tmp []int64) {

	if end <= start {
		return
	}
	mid := start + (end - start) / 2
	mergeSort(arr, start, mid, tmp)
	mergeSort(arr, mid + 1, end, tmp)
	Merge(arr, start, mid, end, tmp)
 }
func Merge(arr []int64, start int, mid int, end int, tmp []int64)  {
	index, i, j := start, start, mid + 1
	for ; i <= mid && j <= end; index++ {
		if arr[i] <= arr[j] {
			tmp[index] = arr[i]
			i++
		} else {
			tmp[index] = arr[j]
			j++
		}
	}
	for ; i <= mid; index++ {
		tmp[index] = arr[i]
		i++
	}
	for ; j <= end; index++ {
		tmp[index] = arr[j]
		j++
	}

	for i := start; i <= end; i++ {
		arr[i] = tmp[i]
	}
}
