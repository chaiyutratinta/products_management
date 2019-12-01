package constance

//RequestErrors error message for validate body
var RequestErrors = map[string]string{
	"Name.required":     "name is required.",
	"Exp.required":      "expire_date id required.",
	"Amount.min":        "amount at least 1.",
	"Price.min":         "price at least 1 THB.",
	"Category.required": "category not match.",
	"Exp.len":           "expire_date length must be 6.",
}
