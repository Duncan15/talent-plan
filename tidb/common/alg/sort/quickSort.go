package sort

//QuickSort the export method for user
func QuickSort(arr []int64, start int, end int)  {
	quickSort(arr, start, end)
}
func quickSort(arr []int64, start int, end int)  {
	if end <= start {
		return
	}
	f, i, j := start, start, end
	for i < j {
		for arr[f] >= arr[i] && i < end {
			i++
		}
		for arr[f] < arr[j] && j > start {
			j--
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[i], arr[j] = arr[j], arr[i]
	arr[f], arr[j] = arr[j], arr[f]
	quickSort(arr, start, j - 1)
	quickSort(arr, j + 1, end)
}

