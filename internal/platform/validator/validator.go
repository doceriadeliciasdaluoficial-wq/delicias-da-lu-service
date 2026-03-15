package validator

type Validator[T any] interface {
	Validate(T) error
}
