package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func insertionSort(nums []int) {
	for j := 1; j < len(nums); j++ {
		key := nums[j]
		i := j - 1
		for ; i >= 0 && nums[i] > key; i-- {
			nums[i+1] = nums[i]
		}
		nums[i+1] = key
	}
}

func benchmarkInsertionSort(b *testing.B) {
	var s [][]int
	for i := 0; i < b.N; i++ {
		s = append(s, rand.Perm(10))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertionSort(s[i])
	}
}

func main() {
	result := testing.Benchmark(benchmarkInsertionSort)
	fmt.Printf("benchmarkInsertionSort\t%v\n", result)
}
