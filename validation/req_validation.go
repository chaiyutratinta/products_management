package validation

import (
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v9"
)

type Validator interface {
	Body(interface{}, interface{}) map[string]string
}

type errorMessage struct {
	errMessage map[string]string
}

//New for init validator
func New(errMessage map[string]string) Validator {

	return &errorMessage{
		errMessage,
	}
}

func (errMsg *errorMessage) Body(body interface{}, bodyType interface{}) map[string]string {
	validate := validator.New()
	err := validate.Struct(body)
	invalid := make(map[string]string)

	if err == nil {
		return invalid
	}

	structType := reflect.TypeOf(bodyType)
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		actualTag := err.ActualTag()

		if field, ok := structType.FieldByName(fieldName); ok {
			key := field.Tag.Get("json")
			errWihTag := fmt.Sprintf("%s.%s", fieldName, actualTag)
			invalid[key] = errMsg.errMessage[errWihTag]
		}
	}

	return invalid
}
