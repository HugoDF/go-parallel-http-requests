package fakestoreapi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type CartResponseProductItem struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CartResponse struct {
	Id       int                       `json:"id"`
	UserId   int                       `json:"userId"`
	Date     string                    `json:"date"`
	Products []CartResponseProductItem `json:"products"`
}

type ProductRating struct {
	Rate  float32 `json:"rate"`
	Count int     `json:"count"`
}

type ProductResponse struct {
	Id          int           `json:"id"`
	Price       float32       `json:"price"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Image       string        `json:"image"`
	Rating      ProductRating `json:"rating"`
}

func LoadCart(cartId int) CartResponse {
	url := fmt.Sprintf("https://fakestoreapi.com/carts/%d", cartId)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var parsedResp CartResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&parsedResp)
	if decodeErr != nil {
		panic(err)
	}
	return parsedResp
}

func LoadProduct(id int) ProductResponse {
	url := fmt.Sprintf("https://fakestoreapi.com/products/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var parsedResp ProductResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&parsedResp)
	if decodeErr != nil {
		panic(err)
	}
	return parsedResp
}

// Naive implementation of loading of products, it is done in a blocking for loop.
func LoadCartAndProductsSequential(cartId int) (CartResponse, []ProductResponse) {
	start := time.Now()

	cartResponse := LoadCart(cartId)

	productResponses := make([]ProductResponse, 0, len(cartResponse.Products))

	for _, product := range cartResponse.Products {
		productRes := LoadProduct(product.ProductId)
		productResponses = append(productResponses, productRes)
	}

	end := time.Now()
	duration := end.Sub(start)
	slog.Info("LoadCartAndProductsSequential runtime",
		"duration", duration,
		"cartId", cartResponse.Id,
		"len(products)", len(productResponses),
	)
	return cartResponse, productResponses
}

// Load products in parallel, since we know the number of calls, we read from the channel that number of times.
func LoadCartAndProductsExhaustChannel(cartId int) (CartResponse, []ProductResponse) {
	start := time.Now()

	cartResponse := LoadCart(cartId)

	productResponsesCh := make(chan ProductResponse, len(cartResponse.Products))
	for _, product := range cartResponse.Products {
		go func(product CartResponseProductItem) {
			productRes := LoadProduct(product.ProductId)
			productResponsesCh <- productRes
		}(product)
	}

	productResponses := make([]ProductResponse, 0, len(cartResponse.Products))
	for i := 0; i < len(cartResponse.Products); i++ {
		comm := <-productResponsesCh
		productResponses = append(productResponses, comm)
	}

	end := time.Now()
	duration := end.Sub(start)
	slog.Info("LoadCartAndProductsExhaustChannel runtime",
		"duration", duration,
		"cartId", cartResponse.Id,
		"len(products)", len(productResponses),
	)
	return cartResponse, productResponses
}

// Load products in parallel, synchronise using a WaitGroup. This ensures we collect results even if one of the calls fails.
func LoadCartAndProductsWaitGroup(cartId int) (CartResponse, []ProductResponse) {
	start := time.Now()

	cartResponse := LoadCart(cartId)

	var wg sync.WaitGroup
	productResponsesCh := make(chan ProductResponse, len(cartResponse.Products))
	for _, product := range cartResponse.Products {
		wg.Add(1)
		go func(product CartResponseProductItem) {
			defer wg.Done()
			productRes := LoadProduct(product.ProductId)
			productResponsesCh <- productRes
		}(product)
	}

	wg.Wait()
	close(productResponsesCh)

	productResponses := make([]ProductResponse, 0, len(cartResponse.Products))
	for chValue := range productResponsesCh {
		productResponses = append(productResponses, chValue)
	}

	end := time.Now()
	duration := end.Sub(start)
	slog.Info("LoadCartAndProductsWaitGroup runtime",
		"duration", duration,
		"cartId", cartResponse.Id,
		"len(products)", len(productResponses),
	)
	return cartResponse, productResponses
}
