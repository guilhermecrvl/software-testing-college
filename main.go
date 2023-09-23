package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func sumEfficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	result, sum := efficientSumBenchmark(num, "efficientSum")
	fmt.Fprint(w, result)
	fmt.Fprintf(w, "Sum: %d\n", sum)
}

func sumInefficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	result, sum := inefficientSumBenchmark(num, "inefficientSum")
	fmt.Fprint(w, result)
	fmt.Fprintf(w, "Sum: %d\n", sum)
}

func efficientSum(n int) int {
	return (n * (n + 1)) / 2
}

func inefficientSum(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

func efficientSumBenchmark(num int, name string) (string, int) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			efficientSum(num)
		}
	})

	return formatBenchmarkResult(&result, memStats, name), efficientSum(num)
}

func inefficientSumBenchmark(num int, name string) (string, int) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			inefficientSum(num)
		}
	})

	return formatBenchmarkResult(&result, memStats, name), inefficientSum(num)
}

func formatBenchmarkResult(result *testing.BenchmarkResult, memStats runtime.MemStats, name string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Function: %s\n", name))
	builder.WriteString(fmt.Sprintf("N: %d\n", result.N))
	builder.WriteString(fmt.Sprintf("Time: %s\n", result.T.String()))
	builder.WriteString(fmt.Sprintf("Memory Alloc: %s\n", formatBytes(memStats.Alloc)))
	builder.WriteString(fmt.Sprintf("Memory Total Alloc: %s\n", formatBytes(memStats.TotalAlloc)))
	builder.WriteString(fmt.Sprintf("Memory Heap Alloc: %s\n", formatBytes(memStats.HeapAlloc)))
	return builder.String()
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div := bytes / unit
	mod := bytes % unit
	return fmt.Sprintf("%.1f KiB", float64(div)+float64(mod)/float64(unit))
}

func main() {
	r := chi.NewRouter()

	r.Get("/sum/efficient/{num:[0-9]+}", sumEfficientHandler)
	r.Get("/sum/inefficient/{num:[0-9]+}", sumInefficientHandler)

	fmt.Println("HTTP server started on :8080")
	http.ListenAndServe(":8080", r)
}
