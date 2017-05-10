package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	log "github.com/sirupsen/logrus"
)

const (
	URL = "https://backend-challenge-fall-2017.herokuapp.com/orders.json"
)

func getResponse(page int) (Response, error) {

	var response Response
	client := http.DefaultClient

	resp, err := client.Get(fmt.Sprintf("%s?page=%d", URL, page))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"url":   URL,
		}).Error("main/getResponse : Failed to http GET")
		return response, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("main/getResponse : Failed to decode response")
		return response, err
	}

	return response, nil
}

func main() {

	var remainingCookies int
	var orders []Order
	var numPages int

	for page := 1; ; page += 1 {
		response, err := getResponse(page)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"page":  page,
			}).Fatal("main/main : Failed to http GET")
		}

		if page == 1 {
			remainingCookies = response.AvailableCookies

			numPages = response.Pagination.Total/response.Pagination.PerPage
			if response.Pagination.Total % response.Pagination.PerPage != 0 {
				numPages += 1
			}
		}

		orders = append(orders, response.Orders...)

		if page == numPages {
			break
		}
	}

	log.WithFields(log.Fields{
		"remainingCookies": remainingCookies,
		"initial orders": orders,
		"len": len(orders),
		"cap": cap(orders),
	}).Print("1. All orders from the paginated API")

	var unfulfilledOrders []Order
	for _, order := range orders {
		if !order.Fulfilled {
			for _, product := range order.Products {
				if product.Title == "Cookie" {
					order.NumCookies = product.Amount
					unfulfilledOrders = append(unfulfilledOrders, order)
					break
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"remainingCookies": remainingCookies,
		"unfulfilledOrders": unfulfilledOrders,
		"len": len(unfulfilledOrders),
		"cap": cap(unfulfilledOrders),
	}).Print("2. Unfulfilled orders")

	sort.Sort(byHighestNumCookiesThenLowestID(unfulfilledOrders))

	log.WithFields(log.Fields{
		"remainingCookies": remainingCookies,
		"unfulfilledOrders": unfulfilledOrders,
		"len": len(unfulfilledOrders),
		"cap": cap(unfulfilledOrders),
	}).Print("3. Sorted unfulfilled orders")

	var unfulfilledOrderIDs []int
	for _, order := range unfulfilledOrders {
		if order.NumCookies > remainingCookies {
			unfulfilledOrderIDs = append(unfulfilledOrderIDs, order.ID)
			continue
		}
		remainingCookies -= order.NumCookies
	}
	sort.Ints(unfulfilledOrderIDs)

	output := Output{
		RemainingCookies: remainingCookies,
		UnfulfilledOrders: unfulfilledOrderIDs,
	}

	answer, err := json.Marshal(output)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"output": output,
		}).Fatal("main/main : Failed to get JSON encoding of output")
	}
	log.WithFields(log.Fields{
		"answer": string(answer),
	}).Print("4. Answer Output")
}
