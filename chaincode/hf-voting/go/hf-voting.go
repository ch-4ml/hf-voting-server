package main

import(
	"encoding/json"
	"fmt"
	"strconv" 		// 문자열 숫자 변환
	"strings" 		// 문자열 포함 검사
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {}

// 투표 클래스
type Vote struct {
	ObjectType	string `json:"docType"`
	ID			string `json:"id"`
	Name		string `json:"name"`
}

// 후보자 클래스
type Candidate struct {
	ObjectType	string `json:"docType"`
	ID			string `json:"id"`
	Name		string `json:"name"`
	VoteID		string `json:"voteID"`	// 투표 ID
	Votes		string `json:"votes"`	// 득표 수
}

// 유권자 클래스
type Electorate struct {
	ObjectType	string `json:"docType"`	 
	ID			string `json:"id"`		// Wallet Identity
	VoteID		string `json:"voteID"`	// 투표 ID
	Voted		string `json:"voted"`	// 투표 여부
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "setVote" {
		return s.setVote(APIstub, args)
	} else if function == "getAllVotes" {
		return s.getAllVotes(APIstub)
	} else if function == "setCandidate" {
		return s.setCandidate(APIstub, args)
	} else if function == "getCandidate" {
		return s.getCandidate(APIstub, args)
	} else if function == "setElectorate" {
		return s.setElectorate(APIstub, args)
	} else if function == "getElecotrate" {
		return s.getElecotrate(APIstub, args)
	} else if function == "vote" {
		return s.vote(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}