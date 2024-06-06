package helpers

func IsCustomerCheck(isCustomer bool, customerID string) string {
	if isCustomer {
		return customerID
	}
	return ""
}