package main

import(
	"encoding/json"
	"fmt"
	"strconv" 		// 문자열 숫자 변환
	"bytes"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {}

// 투표 클래스
type Vote struct {						// 확장 시 총 후보자 수, 총 유권자 수 추가
	ObjectType	string `json:"docType"`
	ID			string `json:"id"`		// 투표 ID
	Name		string `json:"name"`	// 투표 이름(제목)
}

// 후보자 클래스
type Candidate struct {
	ObjectType	string `json:"docType"`
	ID			string `json:"id"`		// 후보자 ID
	Name		string `json:"name"`	// 후보자 이름
	VoteID		string `json:"voteID"`	// 투표 ID
	Votes		string `json:"votes"`	// 득표 수
}

// 유권자 클래스
type Electorate struct {				// 확장 시 기기 ID값의 해시 값을 추가. (분실 / 중복 대안)
	ObjectType	string `json:"docType"`	 
	ID			string `json:"id"`		// Wallet Identity
	VoteID		string `json:"voteID"`	// 투표 ID
	Voted		string `json:"voted"`	// 유권자가 투표한 후보자 이름
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	var result string
	var err error

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "setVote" { 								// 투표 등록: 2 args (ID, Name)
		result, err = setVote(APIstub, args)
	} else if function == "getVote" { 						// 투표 조회: 1 arg (ID)
		result, err = getVote(APIstub, args)
	} else if function == "getVoteByName" {					// 투표 이름으로 조회: 1 arg (Name)
		result, err = getVoteByName(APIstub, args)			
	} else if function == "getAllVotes" { 					// 투표 전체 조회
		result, err = getAllVotes(APIstub)		
	} else if function == "setCandidate" { 					// 후보자 등록: 3 args (ID, Name, VoteID)
		result, err = setCandidate(APIstub, args)		
	} else if function == "getCandidate" {					// 후보자 조회: 1 arg(ID)
		result, err = getCandidate(APIstub, args)		
	} else if function == "getCandidateByNameAndVoteID" { 	// 후보자 이름, 투표 번호로 조회: 2 args (Name, VoteID)
		result, err = getCandidateByNameAndVoteID(APIstub, args)
	} else if function == "getCandidatesByVoteID" { 		// 후보자 전체 조회(해당 투표): 1 arg (VoteID)
		result, err = getCandidatesByVoteID(APIstub, args)		
	} else if function == "setElectorate" { 				// 유권자 등록: 2 args (ID, VoteID)
		result, err = setElectorate(APIstub, args)		
	} else if function == "getElectorate" {					// 유권자 조회(본인만 조회 가능): 1 arg (ID)
		result, err = getElectorate(APIstub, args)		
	} else if function == "isVoted" {						// 투표 여부 검사: 1 arg (ID)
		result, err = isVoted(APIstub, args)
	} else if function == "vote" {							// 투표(후보자, 유권자 수정): 3 args (ID, VoteID, Candidate Name)
		result, err = vote(APIstub, args)		
	} else if function == "getHistoryByVoteID" {			// 투표 기록 조회: 1 arg (VoteID)
		result, err = getHistoryByVoteID(APIstub, args)
	} else { // 현재 투표한 유권자 수 count 등.. 추가
		return shim.Error("Invalid Smart Contract function name.")
	}
	
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

// 투표 등록
func setVote(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 { // ID, Name
		return "", fmt.Errorf("Incorrect arguments. Please input 2 args.")
	}

	// 투표 중복 검사
	getVoteResult, _ := APIstub.GetState(args[0])
	if getVoteResult != nil {
		return "", fmt.Errorf("Already exist: Vote ID")
	}

	// 투표 이름 중복 검사
	getVoteByNameResult, _ := getVoteByName(APIstub, []string{ args[1] })
	if getVoteByNameResult != "" {
		return "", fmt.Errorf("Already exist: Vote name")
	}

	var vote = Vote {
		ObjectType: "Vote",
		ID		  : args[0],
		Name	  :	args[1],
	}

	voteAsBytes, _ := json.Marshal(vote)

	err := APIstub.PutState(args[0], voteAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set vote: %s", err)
	}

	// Composite Key 생성: vote 클래스의 objectType으로 편하게 vote를 조회하기 위해서 사용
	voteIDIndexKey, err := APIstub.CreateCompositeKey("vote~id", []string{ vote.ObjectType, vote.ID })
	if err != nil {
		return "", fmt.Errorf("Failed to set composite key 'vote~id': %s", err)
	}

	// Composite Key 생성: vote 클래스의 name으로 편하게 vote를 조회하기 위해서 사용
	voteNameIDIndexKey, err := APIstub.CreateCompositeKey("votename~id", []string{ vote.Name, vote.ID })
	if err != nil {
		return "", fmt.Errorf("Failed to set composite key 'votename~id': %s", err)
	}

	// 비어있는 byte 배열 생성
	value := []byte{ 0x00 }

	APIstub.PutState(voteIDIndexKey, value)
	APIstub.PutState(voteNameIDIndexKey, value)
	
	return string(voteAsBytes), nil
}

func getVote(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key.")
	}

	id := args[0]

	voteAsBytes, err := APIstub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Vote does not exist: %s", err)
	}

	return string(voteAsBytes), nil
}

func getVoteByName(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}
	name := args[0]
	queriedIDByVoteNameIterator, err := APIstub.GetStateByPartialCompositeKey("votename~id", []string{ name })
	if err != nil {
		return "", fmt.Errorf("Failed to get state by partial composite key 'votename~id': %s", err)
	}
	defer queriedIDByVoteNameIterator.Close()

	var result []byte
	var i int
	for i = 0; queriedIDByVoteNameIterator.HasNext(); i++ {
		res, err := queriedIDByVoteNameIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Error querying 'id' with 'vote name': %s", err)
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(res.Key)
		if err != nil {
			return "", fmt.Errorf("Error splitting composite key('id', 'votename'): %s", err)
		}
		
		returnedName := compositeKeyParts[0]
		returnedID := compositeKeyParts[1]
		fmt.Printf("Found a key from index: %s Name: %s, ID: %s", objectType, returnedName, returnedID)

		voteAsBytes, err := APIstub.GetState(returnedID)
		if err != nil {
			return "", fmt.Errorf("Vote does not exist: %s", err)
		}
		
		result = voteAsBytes
	}

	return string(result), nil
}

func getAllVotes(APIstub shim.ChaincodeStubInterface) (string, error) {
	queriedIDByVoteIterator, err := APIstub.GetStateByPartialCompositeKey("vote~id", []string{ "Vote" })
	if err != nil {
		return "", fmt.Errorf("Failed to get state by partial composite key 'vote~id': %s", err)
	}
	defer queriedIDByVoteIterator.Close()

	var buffer string
	buffer = "["
	comma := false

	var i int
	for i = 0; queriedIDByVoteIterator.HasNext(); i++ {
		res, err := queriedIDByVoteIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Error querying 'id' with 'vote': %s", err)
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(res.Key)
		if err != nil {
			return "", fmt.Errorf("Error splitting composite key('id', 'vote'): %s", err)
		}

		returnedObjectType := compositeKeyParts[0]
		returnedID := compositeKeyParts[1]
		fmt.Printf("Found a key from index: %s ObjectType: %s, Id: %s", objectType, returnedObjectType, returnedID)
		
		if comma == true {
			buffer += ", "
		}
		
		result, err := getVote(APIstub, []string{ returnedID })
		if err != nil {
			return "", fmt.Errorf("Vote does not exist: %s", err)
		}
		buffer += result
		comma = true
	}
	buffer += "]"

	return string(buffer), nil
}

func setCandidate(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 { // ID, Name, VoteID
		return "", fmt.Errorf("Incorrect arguments. Please input 3 args.")
	}

	// 후보자 중복 검사
	getCandidateResult, _ := APIstub.GetState(args[0])
	if getCandidateResult != nil {
		return "", fmt.Errorf("Already exist: Candidate ID")
	}

	// 후보자 이름 중복 검사
	getCandidateByNameAndVoteIDResult, _ := getCandidateByNameAndVoteID(APIstub, []string{ args[1], args[2] })
	if getCandidateByNameAndVoteIDResult != "" {
		return "", fmt.Errorf("Already exist: Candidate name")
	}

	var candidate = Candidate {
		ObjectType: "Candidate",
		ID		  : args[0],
		Name	  :	args[1],
		VoteID	  : args[2],
		Votes	  : "0",
	}

	candidateAsBytes, _ := json.Marshal(candidate)

	err := APIstub.PutState(args[0], candidateAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set candidate: %s", err)
	}

	// Composite Key 생성: candidate 클래스의 objectType과 voteID로 편하게 candidate를 조회하기 위해서 사용
	candidateVoteIDIndexKey, err := APIstub.CreateCompositeKey("candidate~voteid~id", []string{ candidate.ObjectType, candidate.VoteID, candidate.ID })
	if err != nil {
		return "", fmt.Errorf("Failed to set composite key 'candidate~voteid~id': %s", err)
	}

	// Composite Key 생성: candidate 클래스의 name과 voteID로 편하게 candidate를 조회하기 위해서 사용
	candidateNameVoteIDIndexKey, err := APIstub.CreateCompositeKey("candidatename~voteid~id", []string{ candidate.Name, candidate.VoteID, candidate.ID })
	if err != nil {
		return "", fmt.Errorf("Failed to set composite key 'candidatename~voteid~id': %s", err)
	}

	// 비어있는 byte 배열 생성
	value := []byte{ 0x00 }

	APIstub.PutState(candidateVoteIDIndexKey, value)
	APIstub.PutState(candidateNameVoteIDIndexKey, value)

	return string(candidateAsBytes), nil
}

func getCandidate(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key.")
	}

	id := args[0]

	candidateAsBytes, err := APIstub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Candidate does not exist: %s", err)
	}

	return string(candidateAsBytes), nil
}

func getCandidateByNameAndVoteID(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Please input 2 args.")
	}
	name := args[0]
	voteID := args[1]
	queriedCandidateIDByNameAndVoteIDIterator, err := APIstub.GetStateByPartialCompositeKey("candidatename~voteid~id", []string{ name, voteID })
	if err != nil {
		return "", fmt.Errorf("Failed to get state by partial composite key 'candidatename~voteid~id': %s", err)
	}
	defer queriedCandidateIDByNameAndVoteIDIterator.Close()

	var result []byte
	var i int
	for i = 0; queriedCandidateIDByNameAndVoteIDIterator.HasNext(); i++ {
		res, err := queriedCandidateIDByNameAndVoteIDIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Error querying 'id' with 'candidate name': %s", err)
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(res.Key)
		if err != nil {
			return "", fmt.Errorf("Error splitting composite key('id', 'candidatename', 'voteid'): %s", err)
		}
		
		returnedName := compositeKeyParts[0]
		returnedVoteID := compositeKeyParts[1]
		returnedID := compositeKeyParts[2]
		fmt.Printf("Found a key from index: %s Name: %s, VoteID: %s, ID: %s", objectType, returnedName, returnedVoteID, returnedID)

		candidateAsBytes, err := APIstub.GetState(returnedID)
		if err != nil {
			return "", fmt.Errorf("Candidate does not exist: %s", err)
		}
		
		result = candidateAsBytes
	}

	return string(result), nil
}

func getCandidatesByVoteID(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 { // VoteID
		return "", fmt.Errorf("Incorrect arguments. Expecting a key.")
	}

	queriedIDByCandidateIterator, err := APIstub.GetStateByPartialCompositeKey("candidate~voteid~id", []string{ "Candidate", args[0] })
	if err != nil {
		return "", fmt.Errorf("Failed to get state by partial composite key 'candidate~voteid~id': %s", err)
	}
	defer queriedIDByCandidateIterator.Close()

	var buffer string
	buffer = "["
	comma := false

	var i int
	for i = 0; queriedIDByCandidateIterator.HasNext(); i++ {
		res, err := queriedIDByCandidateIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Error querying 'id' with 'candidate': %s", err)
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(res.Key)
		if err != nil {
			return "", fmt.Errorf("Error splitting composite key('id', 'candidate'): %s", err)
		}

		returnedObjectType := compositeKeyParts[0]
		returnedVoteID := compositeKeyParts[1]
		returnedID := compositeKeyParts[2]
		fmt.Printf("Found a key from index: %s ObjectType: %s, VoteID: %s, ID: %s", objectType, returnedObjectType, returnedVoteID, returnedID)
		
		if comma == true {
			buffer += ", "
		}
		
		result, err := getCandidate(APIstub, []string{ returnedID })
		if err != nil {
			return "", fmt.Errorf("Candidate does not exist: %s", err)
		}
		buffer += result
		comma = true
	}
	buffer += "]"

	return string(buffer), nil
}

func setElectorate(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 { // ID, VoteID
		return "", fmt.Errorf("Incorrect arguments. Please input 2 args.")
	}

	// 유권자 중복 검사
	getElectorateResult, _ := APIstub.GetState(args[0])
	if getElectorateResult != nil {
		return "", fmt.Errorf("Already exist: Electorate ID")
	}

	var electorate = Electorate {
		ObjectType: "Electorate",
		ID		  : args[0],
		VoteID	  : args[1],
		Voted	  : "",
	}

	electorateAsBytes, _ := json.Marshal(electorate)

	err := APIstub.PutState(args[0], electorateAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set electorate: %s", err)
	}

	return string(electorateAsBytes), nil
}

func getElectorate(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key.")
	}

	id := args[0]

	electorateAsBytes, err := APIstub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Electorate does not exist: %s", err)
	}

	return string(electorateAsBytes), nil
}

func isVoted(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 { // Electorate ID
		return "", fmt.Errorf("Incorrect arguments. Please input 1 args.")
	}

	id := args[0]

	electorateAsBytes, err := APIstub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Electorate does not exist: %s", err)
	}

	electorateToTransfer := Electorate{}
	err = json.Unmarshal(electorateAsBytes, &electorateToTransfer)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling: %s", err)
	}

	if electorateToTransfer.Voted != "" {
		return "1", nil
	}

	return "0", nil
}

func vote(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 { // ID, VoteID, Candidate Name
		return "", fmt.Errorf("Incorrect arguments. Please input 3 args.")
	}

	id := args[0]
	voteID := args[1]
	candidateName := args[2]

	// 유권자 수정
	electorateAsBytes, err := APIstub.GetState(id)
	if err != nil {
		return "", fmt.Errorf("Electorate does not exist: %s", err)
	}

	electorateToTransfer := Electorate{}
	err = json.Unmarshal(electorateAsBytes, &electorateToTransfer)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling: %s", err)
	}

	electorateToTransfer.Voted = candidateName
	electorateJSONasBytes, _ := json.Marshal(electorateToTransfer)
	err = APIstub.PutState(id, electorateJSONasBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set electorate: %s", err)
	}

	// 후보자 수정
	candidateAsString, err := getCandidateByNameAndVoteID(APIstub, []string{ candidateName, voteID })
	if err != nil {
		return "", fmt.Errorf("Candidate does not exist: %s", err)
	}

	candidateAsBytes := []byte(candidateAsString)
	candidateToTransfer := Candidate{}
	err = json.Unmarshal(candidateAsBytes, &candidateToTransfer)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling: %s", err)
	}

	value := candidateToTransfer.Votes
	iValue, _ := strconv.Atoi(value);
	iValue += 1
	candidateToTransfer.Votes = strconv.Itoa(iValue)

	candidateJSONasBytes, _ := json.Marshal(candidateToTransfer)
	err = APIstub.PutState(candidateToTransfer.ID, candidateJSONasBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set candidate: %s", err)
	}

	return string("Vote succeed"), nil
}

func getHistoryByVoteID(stub shim.ChaincodeStubInterface, args[] string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1 arg")
	}

	id := args[0]

	historyIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	defer historyIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	comma := false
	for historyIterator.HasNext() {
		response, err := historyIterator.Next()
		if err != nil {
			return "", fmt.Errorf("An error occurred while querying the history: %s", err)
		}
		if comma == true {
			buffer.WriteString(", ")
		}
		buffer.WriteString("{\"TxID\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		comma = true
	}
	buffer.WriteString("]")

	return buffer.String(), nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
