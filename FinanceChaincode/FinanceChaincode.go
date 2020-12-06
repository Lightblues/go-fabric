package main

import (

	// "github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type FinanceChaincode struct {
}

//用户
type User struct {
	Name       string   `json:"name"`
	Uid        string   `json:"uid"`
	CompactIDs []string `json:"compactIDs"`
}

//合同
type Compact struct {
	TimeDtamp        int64  `json:"timestamp"`
	Uid              string `json:"uid"`
	LoanAmount       string `json:"loanAmount"`
	ApplyDate        string `json:"applyDate"`
	CompactStartDate string `json:"compactStartDate"`
	CompactEndDate   string `josn:"compactEndDate"`
	ID               string `json:"id"`
}

func (t *FinanceChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	if len(args) != 0 {
		return shim.Error("Parameter error while Init")
	}
	return shim.Success(nil)
}

func (t *FinanceChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functionName, args := stub.GetFunctionAndParameters()
	switch functionName {
	// case "userRegister":
	// 	return userRegister(stub, args)
	case "loan":
		return loan(stub, args)
	case "queryCompact":
		return queryCompact(stub, args)
	case "queryUser":
		return queryUser(stub, args)
	case "enroll":
		return enroll(stub, args)
	case "inquiry":
		return inquiry(stub, args)
	default:
		return shim.Error("Incalid Smart Contract function name.")
	}
}

// 用户注册
// func userRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	if len(args) != 2 {
// 		return shim.Error("Note enough args")
// 	}
// 	name := args[0]
// 	id := args[1]
// 	if name == "" || id == "" {
// 		return shim.Error("Invalid args")
// 	}
// 	if userBytes, err := stub.GetState(id); err != nil || len(userBytes) != 0 {
// 		return shim.Error("User already exist")
// 	}
// 	var user = User{Name: name, Uid: id}
// 	userBytes, err := json.Marshal(user)
// 	if err != nil {
// 		return shim.Error(fmt.Sprint("marshal user error % s", err))
// 	}
// 	fmt.Println()
// 	err = stub.PutState(id, userBytes)
// 	if err != nil {
// 		return shim.Error(fmt.Sprint("put user error % s", err))
// 	}
// 	return shim.Success(nil)
// }

// 登记贷款电子合同
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var compact Compact
	compact.Uid = args[0]
	compact.LoanAmount = args[1]
	compact.ApplyDate = args[2]
	compact.CompactStartDate = args[3]
	compact.CompactEndDate = args[4]
	compact.TimeDtamp = time.Now().Unix()
	compact.ID = args[5]

	ownBytes, err := stub.GetState(compact.Uid)
	if err != nil || len(ownBytes) == 0 {
		return shim.Error("user not found")
	}
	compactBytes, err := json.Marshal(&compact) // json 序列化
	if err != nil {
		return shim.Error("Json serialize Compact fail while Loan")
	}
	if compactBytes, err := stub.GetState(args[5]); err != nil || len(compactBytes) != 0 {
		return shim.Error("Compact already exist")
	}
	err = stub.PutState(args[5], compactBytes)
	if err != nil {
		return shim.Error(fmt.Sprint("Put Compact error % s", err))
	}

	owner := new(User)
	if err := json.Unmarshal(ownBytes, owner); err != nil {
		return shim.Error(fmt.Sprint("unmarshal user error % s", err))
	}
	owner.CompactIDs = append(owner.CompactIDs, compact.ID)
	ownerBytes, err := json.Marshal(owner)

	err = stub.PutState(compact.Uid, ownerBytes)

	return shim.Success([]byte("记录贷款成功"))

}

// 查询贷款电子合同
func queryCompact(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of argumnets. Expecting 1")
	}
	compactID := args[0]
	compactBytes, err := stub.GetState(compactID)
	if err != nil || len(compactBytes) == 0 {
		return shim.Error("compact not found")
	}
	return shim.Success(compactBytes)
}

// 查询用户
func queryUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	userID := args[0]
	userBytes, err := stub.GetState(userID)
	if err != nil || len(userBytes) == 0 {
		return shim.Error("compact not found")
	}
	return shim.Success(userBytes)
}

/* 下面是作业要求实现的 */
// 用户注册，以下直接修改自课本
// peer chaincode invoke -C mychannel -n mychannel -c '{"Args":["enroll", "user1","agagahagahwrghag"]}'
func enroll(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Note enough args")
	}
	name := args[0]
	id := args[1]
	if name == "" || id == "" {
		return shim.Error("Invalid args")
	}
	if userBytes, err := stub.GetState(id); err != nil || len(userBytes) != 0 {
		return shim.Error("User already exist")
	}
	var user = User{Name: name, Uid: id}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprint("marshal user error % s", err))
	}
	err = stub.PutState(id, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprint("put user error % s", err))
	}
	return shim.Success(nil)
}

// 查询贷款电子合同，注意这里的两个参数分别是用户 ID 和合同 ID，因此查询逻辑是，先基于合同 ID 查询对应的合同，再利用用户 ID 验证该合同是否属于该用户
//peer chaincode invoke -C mychannel -n mychannel -c '{"Args":["inquiry", "agagahagahwrghag","000001"]}'
func inquiry(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Args error, need 2")
	}
	userID := args[0]
	compactID := args[1]

	compactBytes, err := stub.GetState(compactID)
	if err != nil || len(compactBytes) == 0 {
		return shim.Error("Cannot find compact")
	}
	compact := new(Compact)
	if err := json.Unmarshal(compactBytes, compact); err != nil {
		return shim.Error("Json unmarshal error")
	}
	if userID != compact.Uid {
		return shim.Error(fmt.Sprint("The compact % s does not belong to % s !", compactID, userID))
	}
	return shim.Success(compactBytes)
}

func main() {
	if err := shim.Start(new(FinanceChaincode)); err != nil {
		fmt.Printf("Error createing new Smart Contract: %s", err)
	}
}
