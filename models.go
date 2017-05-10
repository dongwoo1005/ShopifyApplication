package main

type Response struct {
	AvailableCookies int     `json:"available_cookies"`
	Orders           []Order `json:"orders"`
	Pagination       `json:"pagination"`
}

type Order struct {
	ID            int       `json:"id"`
	Fulfilled     bool      `json:"fulfilled"`
	CustomerEmail string    `json:"customer_email"`
	Products      []Product `json:"products"`
	NumCookies    int       `json:"num_cookies"`
}

type Product struct {
	Title     string  `json:"title"`
	Amount    int     `json:"amount"`
	UnitPrice float64 `json:"unit_price"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

type Output struct {
	RemainingCookies  int   `json:"remaining_cookies"`
	UnfulfilledOrders []int `json:"unfulfilled_orders"`
}

type byHighestNumCookiesThenLowestID []Order

func (p byHighestNumCookiesThenLowestID) Len() int {
	return len(p)
}

func (p byHighestNumCookiesThenLowestID) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// sort unfulfilledOrders by higher numCookies, then lower ID
func (p byHighestNumCookiesThenLowestID) Less(i, j int) bool {
	if p[i].NumCookies > p[j].NumCookies {
		return true
	}
	if p[i].ID < p[j].ID {
		return true
	}
	return false
}
