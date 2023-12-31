package fakestoreapi_test

import (
	"encoding/json"
	"go-parallel-http-requests/fakestoreapi"
	"strings"
	"testing"

	"github.com/h2non/gock"
)

const mockCartResponse = `{
	"id": 1,
	"userId": 1,
	"date": "2020-03-02T00:00:00.000Z",
	"products": [
		{
			"productId": 1,
			"quantity": 4
		},
		{
			"productId": 2,
			"quantity": 1
		}
	],
	"__v": 0
}`

const mockProduct1Response = `{
  "id": 1,
  "title": "Fjallraven - Foldsack No. 1 Backpack, Fits 15 Laptops",
  "price": 109.95,
  "description": "Your perfect pack for everyday use and walks in the forest. Stash your laptop (up to 15 inches) in the padded sleeve, your everyday",
  "category": "men's clothing",
  "image": "https://fakestoreapi.com/img/81fPKd-2AYL._AC_SL1500_.jpg",
  "rating": {
    "rate": 3.9,
    "count": 120
  }
}
`
const mockProduct2Response = `{
  "id": 2,
  "title": "Mens Casual Premium Slim Fit T-Shirts ",
  "price": 22.3,
  "description": "Slim-fitting style, contrast raglan long sleeve, three-button henley placket, light weight & soft fabric for breathable and comfortable wearing. And Solid stitched shirts with round neck made for durability and a great fit for casual fashion wear and diehard baseball fans. The Henley style round neckline includes a three-button placket.",
  "category": "men's clothing",
  "image": "https://fakestoreapi.com/img/71-3HjGNDUL._AC_SY879._SX._UX._SY._UY_.jpg",
  "rating": {
    "rate": 4.1,
    "count": 259
  }
}`

func MockRes(mock string) any {
	var res (any)
	json.Unmarshal([]byte(mock), &res)
	return res
}

func TestSequential(t *testing.T) {
	gock.New("https://fakestoreapi.com").
		Get("carts/1").Reply(200).JSON(MockRes(mockCartResponse))
	gock.New("https://fakestoreapi.com").
		Get("products/1").Reply(200).JSON(MockRes(mockProduct1Response))
	gock.New("https://fakestoreapi.com").
		Get("products/2").Reply(200).JSON(MockRes(mockProduct2Response))
	defer gock.Off()

	cart, products := fakestoreapi.LoadCartAndProductsSequential(1)
	if !strings.Contains(cart.Date, "2020-03") {
		t.Fatalf(`cart.Date: %s, expected: %s`, cart.Date, "2020-03*")
	}

	if len(products) != 2 {
		t.Fatalf(`products expected: %v, actual: %v`, 2, len(products))
	}
}

func TestExhaust(t *testing.T) {
	gock.New("https://fakestoreapi.com").
		Get("carts/1").Reply(200).JSON(MockRes(mockCartResponse))
	gock.New("https://fakestoreapi.com").
		Get("products/1").Reply(200).JSON(MockRes(mockProduct1Response))
	gock.New("https://fakestoreapi.com").
		Get("products/2").Reply(200).JSON(MockRes(mockProduct2Response))
	defer gock.Off()

	cart, products := fakestoreapi.LoadCartAndProductsExhaustChannel(1)
	if !strings.Contains(cart.Date, "2020-03") {
		t.Fatalf(`cart.Date: %s, expected: %s`, cart.Date, "2020-03*")
	}

	if len(products) != 2 {
		t.Fatalf(`products expected: %v, actual: %v`, 2, len(products))
	}
}

func TestWaitGroup(t *testing.T) {
	gock.New("https://fakestoreapi.com").
		Get("carts/1").Reply(200).JSON(MockRes(mockCartResponse))
	gock.New("https://fakestoreapi.com").
		Get("products/1").Reply(200).JSON(MockRes(mockProduct1Response))
	gock.New("https://fakestoreapi.com").
		Get("products/2").Reply(200).JSON(MockRes(mockProduct2Response))
	defer gock.Off()

	cart, products := fakestoreapi.LoadCartAndProductsWaitGroup(1)
	if !strings.Contains(cart.Date, "2020-03") {
		t.Fatalf(`cart.Date: %s, expected: %s`, cart.Date, "2020-03*")
	}

	if len(products) != 2 {
		t.Fatalf(`products expected: %v, actual: %v`, 2, len(products))
	}
}
