package main

import (
	"runtime"
	"testing"
)

const sliceSize = 1000000

func BenchmarkSum(b *testing.B) {
	nums := make([]int, sliceSize)
	for i := 0; i < sliceSize; i++ {
		nums[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(nums)
	}
}

func BenchmarkSumConcurrently(b *testing.B) {
	nums := make([]int, sliceSize)
	for i := 0; i < sliceSize; i++ {
		nums[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumConcurrently(nums, runtime.NumCPU())
	}
}

const requestNum = 20

func BenchmarkDo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Do(requestNum)
	}
}

func BenchmarkDoConcurrently(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoConcurrently(requestNum)
	}
}
