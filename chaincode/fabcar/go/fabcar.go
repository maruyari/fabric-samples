package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Student describes basic details of what makes up a car
type Student struct {
	Name  string `json:"name"`
	Year  string `json:"year"`
	Board string `json:"board"`
	Mark  string `json:"mark"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Student
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	students := []Student{
		{Name: "changed", Year: "2018", Board: "CBSE", Mark: "99"},
		{Name: "jane", Year: "2017", Board: "ICSE", Mark: "92"},
		{Name: "Tan", Year: "2018", Board: "CBSE", Mark: "85"},
		{Name: "jon", Year: "2018", Board: "ICSE", Mark: "86"},
		{Name: "Om", Year: "2018", Board: "CBSE", Mark: "89"},
		{Name: "Vaish", Year: "2018", Board: "CBSE", Mark: "94"},
		{Name: "Rut", Year: "2016", Board: "GSB", Mark: "93"},
		{Name: "Rat", Year: "2015", Board: "CBSE", Mark: "84"},
		{Name: "Vir", Year: "2018", Board: "MSB", Mark: "99"},
		{Name: "Jo", Year: "2018", Board: "CBSE", Mark: "99"},
	}

	for i, student := range students {
		studentAsBytes, _ := json.Marshal(student)
		err := ctx.GetStub().PutState("Student"+strconv.Itoa(i), studentAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateMarksheet adds a new car to the world state with given details
func (s *SmartContract) CreateMarksheet(ctx contractapi.TransactionContextInterface, StudNumber string, name string, year string, board string, mark string) error {
	car := Student{
		Name:  name,
		Year:  year,
		Board: board,
		Mark:  mark,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(StudNumber, carAsBytes)
}

// QueryMarksheet returns the car stored in the world state with given id
func (s *SmartContract) QueryMarksheet(ctx contractapi.TransactionContextInterface, studNumber string) (*Student, error) {
	studAsBytes, err := ctx.GetStub().GetState(studNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if studAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", studNumber)
	}

	car := new(Student)
	_ = json.Unmarshal(studAsBytes, car)

	return car, nil
}

// QueryFullMarksheet returns all cars found in world state
func (s *SmartContract) QueryFullMarksheet(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "Student0"
	endKey := "Student99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		student := new(Student)
		_ = json.Unmarshal(queryResponse.Value, student)

		queryResult := QueryResult{Key: queryResponse.Key, Record: student}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeStudentMarks updates the owner field of car with given id in world state
func (s *SmartContract) ChangeStudentMarks(ctx contractapi.TransactionContextInterface, studNumber string, newMark string) error {
	student, err := s.QueryMarksheet(ctx, studNumber)

	if err != nil {
		return err
	}

	student.Mark = newMark

	studAsBytes, _ := json.Marshal(student)

	return ctx.GetStub().PutState(studNumber, studAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
