package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func sumEfficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	result, sum, memStats, cpuUsage := sumBenchmark(num, "efficientSum")
	formattedResult := formatBenchmarkResult("efficientSum", num, result, sum, memStats, cpuUsage)
	fmt.Fprint(w, formattedResult)
}

func sumInefficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	result, sum, memStats, cpuUsage := sumBenchmark(num, "inefficientSum")
	formattedResult := formatBenchmarkResult("inefficientSum", num, result, sum, memStats, cpuUsage)
	fmt.Fprint(w, formattedResult)
}

func sumBenchmark(num int, function string) (time.Duration, int, runtime.MemStats, int) {
	startCPU := runtime.NumCPU()
	startTime := time.Now()
	var memStats runtime.MemStats
	sum := 0

	if function == "efficientSum" {
		sum = efficientSum(num)
	} else {
		sum = inefficientSum(num)
	}

	runtime.ReadMemStats(&memStats)
	elapsedTime := time.Since(startTime)
	endCPU := runtime.NumCPU()
	cpuUsage := endCPU - startCPU

	return elapsedTime, sum, memStats, cpuUsage
}

func efficientSum(n int) int {
	sum := (n * (n + 1)) / 2
	return sum
}

func inefficientSum(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

func formatBenchmarkResult(name string, num int, elapsedTime time.Duration, sum int, memStats runtime.MemStats, cpuUsage int) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Function: %s\n", name))
	builder.WriteString(fmt.Sprintf("N: %d\n", num))
	builder.WriteString(fmt.Sprintf("Sum: %d\n", sum))
	builder.WriteString(fmt.Sprintf("Time: %s\n", elapsedTime))
	builder.WriteString(fmt.Sprintf("Memory Alloc: %s\n", formatBytes(memStats.Alloc)))
	builder.WriteString(fmt.Sprintf("Memory Total Alloc: %s\n", formatBytes(memStats.TotalAlloc)))
	builder.WriteString(fmt.Sprintf("Memory Heap Alloc: %s\n", formatBytes(memStats.HeapAlloc)))
	builder.WriteString(fmt.Sprintf("CPU Usage: %d\n", cpuUsage))
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
	r.Get("/sum/efficient/{num}", sumEfficientHandler)
	r.Get("/sum/inefficient/{num:[0-9]+}", sumInefficientHandler)

	fmt.Println("HTTP server started on :8080")
	http.ListenAndServe(":8080", r)
}
