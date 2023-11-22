package option

// A Func can be used to mutate a value, generally at creation.
type Func[T any] func(T) error

// Apply will apply all option functions given to the value provided,
// and return the first error it encounters (if any). If an error is
// encountered, the remaining option functions will not be applied.
func Apply[T any](value T, opts ...Func[T]) error {
	for _, f := range opts {
		err := f(value)
		if err != nil {
			return err
		}
	}

	return nil
}
