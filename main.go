package main

import (
	"go-parallel-http-requests/fakestoreapi"
)

func main() {
	fakestoreapi.LoadCartAndProductsSequential(1)
	fakestoreapi.LoadCartAndProductsExhaustChannel(1)
	fakestoreapi.LoadCartAndProductsWaitGroup(1)
}
