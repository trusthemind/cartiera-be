package helpers

import (
	"sort"

	"github.com/stripe/stripe-go"
)

func SortByTime(array []*stripe.PaymentMethod) []*stripe.PaymentMethod {
	var sorted []*stripe.PaymentMethod

	for _, v := range array {
		sorted = append(sorted, v)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Created > sorted[j].Created
	})

	return sorted
}
