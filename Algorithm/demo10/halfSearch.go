package main

import "fmt"

func main() {
	a := []int{19, 2, 38, 41, 60, 16, 27, 58}
	fmt.Println(a)
	QuickSort(a, 0, len(a)-1)
	fmt.Print(a)
}

func partion(a []int, low, high int) int {
	temp := a[low]
	for low < high {
		for low < high && a[high] > temp {
			high--
		}
		a[low] = a[high]

		for low < high && a[low] < temp {
			low++
		}
		a[high] = a[low]
	}
	a[low] = temp
	return low
}

func QuickSort(a []int, low, high int) {
	if low < high {
		partion := partion(a, low, high)
		QuickSort(a, low, partion-1)
		QuickSort(a, partion+1, high)
	}
}
