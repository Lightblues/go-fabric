package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

/*下面的主要内容修改自 SimpleChaincode，尝试着对于要求的三个 Invoke 方法进行了简单的测试*/

func performInit(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	//通过stub调用链码中的Init
	res := stub.MockInit("1", args)
	/*- 1为uuid，用于链码开始前和结束后开始事务的标志，无实际意义
	  - args为初始化需要的参数*/
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		//如果没有启动成功，则报告失败
		t.FailNow()
	}
}

//确认指定key的状态值是否为预期的值
func checkState(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	//获取链码维护的指定key的状态值
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	//如果不是期望的值，则报告失败
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func Test_Init(t *testing.T) {
	//构建链码
	cc := new(FinanceChaincode)
	// 将其部署到模拟MockStub平台， 传入名称和链码实体
	stub := shimtest.NewMockStub("FinanceChaincode", cc)
	// //带参数（“a","10)调用链码中的Init，并确认是否成功
	// performInit(t, stub, [][]byte{[]byte("a"), []byte("10")})
	// //确认初始化后，指定key的状态值是否为预期的值，即健“a" 的vlaue 应该为“10”
	// checkState(t, stub, "a", "10")
	performInit(t, stub, [][]byte{})
}

func Test_Enroll(t *testing.T) {
	cc := new(FinanceChaincode)
	// 获取MockStub对象， 传入名称和链码实体
	stub := shimtest.NewMockStub("SimpleChaincode", cc)
	performInit(t, stub, [][]byte{})

	//调用invoke方法中的set方法
	stub.MockInvoke("1", [][]byte{[]byte("enroll"), []byte("user1"), []byte("agagahagahwrghag")}) //注册用户

	user := User{Uid: "agagahagahwrghag", Name: "user1"}
	userBytes, _ := json.Marshal(user)
	checkState(t, stub, "agagahagahwrghag", string(userBytes))

}

func Test_Loan(t *testing.T) {
	cc := new(FinanceChaincode)
	stub := shimtest.NewMockStub("FinanceChaincode", cc)
	performInit(t, stub, [][]byte{})

	stub.MockInvoke("1", [][]byte{[]byte("enroll"), []byte("user1"), []byte("agagahagahwrghag")}) //注册用户
	res1 := stub.MockInvoke("1", [][]byte{[]byte("loan"), []byte("agagahagahwrghag"), []byte("2020"), []byte("20200501090823"), []byte("20200501090823"), []byte("20201001090823"), []byte("000001")})
	if string(res1.Payload) != "记录贷款成功" {
		fmt.Println("Error")
		t.FailNow()
	}

	res := stub.MockInvoke("1", [][]byte{[]byte("inquiry"), []byte("agagahagahwrghag"), []byte("000001")})

	// 查看 MockInvoke 返回结果信息
	fmt.Println(string(res.Payload), res.Status, res.Message) // 参看 peer.Response 的结构定义

	var compact Compact
	json.Unmarshal(res.Payload, &compact)
	fmt.Println(compact)              //
	if compact.LoanAmount != "2020" { // 测试成功不会显示 print 结果；调试过程中可故意出错，从而打印信息
		fmt.Println("Wrong")
		t.FailNow()
	}
}
