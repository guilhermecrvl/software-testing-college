package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func sumEfficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	sum := efficientSum(num)
	formattedResult := formatResult("efficientSum", num, sum)
	fmt.Fprint(w, formattedResult)
}

func sumInefficientHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}
	sum := inefficientSum(num)
	formattedResult := formatResult("inefficientSum", num, sum)
	fmt.Fprint(w, formattedResult)
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

func formatResult(name string, num int, sum int) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Function: %s\n", name))
	builder.WriteString(fmt.Sprintf("N: %d\n", num))
	builder.WriteString(fmt.Sprintf("Sum: %d\n", sum))
	return builder.String()
}

func main() {
	r := chi.NewRouter()
	r.Get("/sum/efficient/{num:[0-9]+}", sumEfficientHandler)
	r.Get("/sum/inefficient/{num:[0-9]+}", sumInefficientHandler)

	fmt.Println("HTTP server started on :3030")
	http.ListenAndServe(":3030", r)
}
