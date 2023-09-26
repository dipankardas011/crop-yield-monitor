package main

import (
	"testing"
)

func TestHash(t *testing.T) {
	s := "1234"
	a := "1234"

	for i := 0; i < 100000; i++ {
		t.Run("salt testing", func(t *testing.T) {

			salt1, err := GenerateRandomString()
			if err != nil {
				panic(err)
			}

			salt2, err := GenerateRandomString()
			if err != nil {
				panic(err)
			}
			if genHash(s+salt1) == genHash(a+salt2) {
				t.Error("match with salts being same")
			}
		})
	}

}
