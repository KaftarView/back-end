package application_math

import (
	"first-project/src/repository"
	"log"
)

type AddService struct {
	userRepository *repository.UserRepository
}

func NewAddService(userRepository *repository.UserRepository) *AddService {
	return &AddService{
		userRepository: userRepository,
	}
}

func (addService *AddService) Add(num1 int, num2 int) int {
	tables := addService.userRepository.Test()
	results := addService.userRepository.Test2()

	for _, table := range tables {
		log.Println("Table:", table)
	}
	for _, record := range results {
		log.Printf("Name: %s, Age: %d\n", record.Name, record.Age)
	}
	return num1 + num2
}
