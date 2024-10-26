package http

import (
	"github.com/SlavaShagalov/ds-lab1/internal/models"
)

type Person struct {
	ID      int
	Name    string
	Age     int
	Address string
	Work    string
}

// API requests
type createRequest struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type partialUpdateRequest struct {
	Name    *string `json:"name"`
	Age     *int    `json:"age"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
}

// API responses
type createResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func newCreateResponse(person *models.Person) *createResponse {
	return &createResponse{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}

type listResponse struct {
	Persons []models.Person `json:"persons"`
}

func newListResponse(persons []models.Person) *listResponse {
	return &listResponse{
		Persons: persons,
	}
}

type getResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func newGetResponse(person *models.Person) *getResponse {
	return &getResponse{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}
