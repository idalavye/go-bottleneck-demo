/*
go build -gcflags="-m" ./pkg/example/example4
*/
package example4

import (
	"math/rand"
)

const (
	userCount = 10_000
)

// User struct simulates a real-life entity
// Simple struct
type User struct {
	ID    int
	Name  string
	Score int
}

// Value semantic factory
func GenerateUsersValue(n int) []User {
	users := make([]User, n)
	for i := range users {
		users[i] = User{
			ID:    i,
			Name:  "User",
			Score: rand.Intn(100),
		}
	}
	return users // Value returned, usually stays on stack if not escaping
}

// Pointer semantic factory
func GenerateUsersPointer(n int) *[]User {
	users := make([]User, n)
	for i := range users {
		users[i] = User{
			ID:    i,
			Name:  "User",
			Score: rand.Intn(100),
		}
	}
	return &users
}

// Sums scores for value slice
func SumScoresValue() int {
	users := make([]User, userCount)
	for i := range users {
		users[i] = User{
			ID:    i,
			Name:  "User",
			Score: rand.Intn(100),
		}
	}
	sum := 0
	for _, u := range users {
		sum += u.Score
	}
	return sum
}

// Sums scores for pointer slice
func SumScoresPointer() int {
	users := GenerateUsersPointer(userCount)
	sum := 0
	for _, u := range *users {
		sum += u.Score
	}
	return sum
}

// SumFixedArray sums scores for a fixed-size array
func SumFixedArray() int {
	var users [userCount]User
	for i := range users {
		users[i] = User{
			ID:    i,
			Name:  "User",
			Score: rand.Intn(100),
		}
	}
	sum := 0
	for _, u := range users {
		sum += u.Score
	}
	return sum
}
