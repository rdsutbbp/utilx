package validatex

import "github.com/go-playground/validator/v10"

// CheckStructTags check ui model of args
func CheckStructTags(obj interface{}) (bool, error) {
	v := validator.New()
	if err := v.Struct(obj); err != nil {
		return false, err
	}
	return true, nil
}
