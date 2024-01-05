package fakestoreapi_test

import (
	"fmt"
	"go-parallel-http-requests/fakestoreapi"
	"time"

	"github.com/h2non/gock"
)

func Example() {
	gock.New("https://fakestoreapi.com").
		Get("carts/1").Reply(200).JSON(MockRes(mockCartResponse)).Delay(338 * time.Millisecond).Mock.Request().Times(3)
	gock.New("https://fakestoreapi.com").
		Get("products/1").Reply(200).JSON(MockRes(mockProduct1Response)).Delay(169 * time.Millisecond).Mock.Request().Times(3)
	gock.New("https://fakestoreapi.com").
		Get("products/2").Reply(200).JSON(MockRes(mockProduct2Response)).Delay(272 * time.Millisecond).Mock.Request().Times(3)
	// product 3 req -> 290ms
	defer gock.Off()

	cart, products := fakestoreapi.LoadCartAndProductsSequential(1)
	fmt.Println(cart.Date, len(products))

	cart, products = fakestoreapi.LoadCartAndProductsExhaustChannel(1)
	fmt.Println(cart.Date, len(products))

	cart, products = fakestoreapi.LoadCartAndProductsWaitGroup(1)
	fmt.Println(cart.Date, len(products))

	// Output: 2020-03-02T00:00:00.000Z 2
	// 2020-03-02T00:00:00.000Z 2
	// 2020-03-02T00:00:00.000Z 2
}
