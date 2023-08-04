package id

import "github.com/google/uuid"

func New() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic("failed to generate run ID: " + err.Error())
	}
	return u.String()
}
