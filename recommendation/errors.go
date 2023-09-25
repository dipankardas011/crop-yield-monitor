package main

type apiError struct {
	Status int
	Err    string
}

func (e apiError) Error() string {
	return e.Err
}
