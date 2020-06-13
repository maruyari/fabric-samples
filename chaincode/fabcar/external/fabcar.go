/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ServerConfig struct {
	CCID    string
	Address string
}

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Student struct {
	Name  string `json:"make"`
	Year  string `json:"model"`
	Board string `json:"colour"`
	Mark  string `json:"owner"`
	Roll  string `json:"rollno"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Student
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Student{
		{Name: "changed", Year: "2018", Board: "CBSE", Mark: "99", Roll: "290319087"},
		{Name: "jane", Year: "2017", Board: "ICSE", Mark: "92", Roll: "290393087"},
		{Name: "Tan", Year: "2018", Board: "CBSE", Mark: "85", Roll: "2903914087"},
		{Name: "jon", Year: "2018", Board: "ICSE", Mark: "86", Roll: "290329087"},
		{Name: "Om", Year: "2018", Board: "CBSE", Mark: "89", Roll: "290379087"},
		{Name: "Vaish", Year: "2018", Board: "CBSE", Mark: "94", Roll: "20039087"},
		{Name: "Rut", Year: "2016", Board: "GSB", Mark: "93", Roll: "29031087"},
		{Name: "Rat", Year: "2015", Board: "CBSE", Mark: "84", Roll: "29035087"},
		{Name: "Vir", Year: "2018", Board: "MSB", Mark: "99", Roll: "29039287"},
		{Name: "Jo", Year: "2018", Board: "CBSE", Mark: "99", Roll: "29039487"},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("Student"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateMarksheet adds a new car to the world state with given details
func (s *SmartContract) CreateMarksheet(ctx contractapi.TransactionContextInterface, StudNumber string, name string, year string, board string, mark string, rollno string) error {
	student := Student{
		Name:  name,
		Year:  year,
		Board: board,

		Mark: mark,
		Roll: rollno,
	}

	StudentAsBytes, _ := json.Marshal(student)

	return ctx.GetStub().PutState(StudNumber, StudentAsBytes)
}

// QueryMarksheet returns the car stored in the world state with given id
func (s *SmartContract) QueryMarksheet(ctx contractapi.TransactionContextInterface, rollno string) (*Student, error) {
	studAsBytes, err := ctx.GetStub().GetState(rollno)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if studAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", rollno)
	}

	student := new(Student)
	_ = json.Unmarshal(studAsBytes, student)

	return student, nil
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

		car := new(Student)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeStudentMarks updates the owner field of car with given id in world state
func (s *SmartContract) ChangeStudentMarks(ctx contractapi.TransactionContextInterface, rollno string, newMark string) error {
	car, err := s.QueryMarksheet(ctx, rollno)

	if err != nil {
		return err
	}

	car.Mark = newMark

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(rollno, carAsBytes)
}

func main() {
	// See chaincode.env.example
	config := ServerConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
