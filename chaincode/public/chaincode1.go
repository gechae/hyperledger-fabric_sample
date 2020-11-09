package main

import (
	"fmt"
	"log"
	"encoding/json"
	"reflect"
	"bytes"
	"strconv"

	s "github.com/hyperledger/fabric/core/chaincode/shim"
	p "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

type struct1 struct {
}
type struct2 struct {
}
type SummaryStruct struct {
	dates					string					`json:"dates"`
	cnt						int					`json:"cnt"`
	amt						int					`json:"amt"`
}
type RespMsg struct {
	Success		bool			`json:"success"`
	Code		int 			`json:"code"`
	Msg			string			`json:"msg"`
}
type SearchStruct3 struct {
}
func (c *Chaincode) Init(st s.ChaincodeStubInterface) p.Response {
	return s.Success(nil)
}

func (c *Chaincode) Invoke(st s.ChaincodeStubInterface) p.Response {
	function, args := st. GetFunctionAndParameters()

	switch function {
	
	case "1": 
		return c.func1(st, args)
	case "2":
		return c.func2(st, args)
	case "2":
		return c.func3(st, args)


	case "3": 
		return c.func4(st, args)
	case "4", "5": 
		return c.func5(st,args)
	case "6","7":
		return c.func6(st, args)
	}

	response := resultMapMsg(false,8888,"Invalid Smart Contract function name.",[]string{})
	return s.Success(response)
}

func (c *Chaincode) func1(st s.ChaincodeStubInterface, args []string) p.Response{
	log.Print("====================== func1 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21000,"Incorrect number of argumnets 1.",""))
	}

	log.Print("arg[0] :" +args[0])
	var data map[string][]struct
	err := json.Unmarshal([]byte(args[0]), &data)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21001,"Filed to decode Json of: "+ string(args[0]),""))
	}

	Data1, exists := data["data"]
	if !exists {
		return s.Success(resultMapMsg(false,21012,"Is not Data1 Data",""))
	}
	
	for i:=0; i < len(Data1); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMapMsg(false,21009,"Can Not Create CompositeKey",err.Error()))
		}


		bytestr, err := json.Marshal(Data1[i])
		if err != nil {
			return s.Success(resultMapMsg(false,21010,"Convert Maershal Error",err.Error()))
		}

		err = st.PutState(compositeKey,bytestr)
		if err != nil {
			return s.Success(resultMapMsg(false,21011,"Can not put state",err.Error()))
		}
	}

	response := resultMapMsg(true,9999,"Success func1 ","")
	return s.Success(response)
}

func (c *Chaincode) func2(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func2 ======================")

	
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21100,"Incorrect number of argumnets 1.",[]string{}))
	}

	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21102,"Filed to decode Json of: "+ string(args[0]),""))
	}

	var summary SummaryStruct
	
	// Rich Query 실행
	queryiterater, err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMapMsg(false,21105,"query operation on private data failed. Error accessing state",[]string{}))
	}
	defer queryiterater.Close()


	for queryiterater.HasNext() {
		response, iterErr := queryiterater.Next()
		if iterErr != nil {
			return s.Success(resultMapMsg(false,21106,"query operation on private data failed. Error accessing state",[]string{}))
		}
		var struct1Info Struct1
		err := json.Unmarshal(response.Value, &struct1Info)
		if err != nil {
			return s.Success(resultMapMsg(false,21107,"Unmarshal Error %s, Convert Marshal",[]string{}))
		}
		amt,_ := strconv.Atoi()
		summary.Mnrc_amt += amt
		summary.Mnrc_cnt += 1
	}
	var summaryList []SummaryStruct
	summaryList = append(summaryList, summary)

	response := resultMapMsg(true,9999,"Success func2",summaryList)
	log.Print(string(response))
	return s.Success(response)
}
func (c *Chaincode) func3(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func3 ======================")

	
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21100,"Incorrect number of argumnets 1.",[]string{}))
	}

	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21101,"Filed to decode Json of: "+ string(args[0]),""))
	}

	// Rich Query 실행
	queryiterater, err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMapMsg(false,21103,"query operation on private data failed. Error accessing state",[]string{}))
	}
	defer queryiterater.Close()

	var struct1InfoList []Struct1
	for queryiterater.HasNext() {
		response, iterErr := queryiterater.Next()
		if iterErr != nil {
			return s.Success(resultMapMsg(false,21104,"query operation on private data failed. Error accessing state",[]string{}))
		}
		// struct로 변화(Unmarshal)
		var struct1Info Struct1
		err := json.Unmarshal(response.Value, &struct1Info)
		if err != nil {
			return s.Success(resultMapMsg(false,21105,"Unmarshal Error %s, Convert Marshal",[]string{}))
		}
		struct1InfoList = append(struct1InfoList,struct1Info)
	}
	
	response := resultMapMsg(true,9999,"Success func3",struct1InfoList)
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func4(st s.ChaincodeStubInterface, args []string) p.Response{
	log.Print("====================== func4 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21200,"Incorrect number of argumnets 1.",""))
	}

	log.Print("Request Argument : " + args[0])
	var data map[string][]Struct2
	err := json.Unmarshal([]byte(args[0]), &data)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21201,"Filed to decode Json of: "+ string(args[0]),""))
	}
	Data2, exists := data["data"]
	if !exists {
		return s.Success(resultMapMsg(false,21012,"Is not refund Data",""))
	}
	
	for i:=0; i < len(Data2); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMapMsg(false,21208,"Can Not Create CompositeKey",err.Error()))
		}

		bytestr, err := json.Marshal(Data2[i])
		if err != nil {
			return s.Success(resultMapMsg(false,21209,"Convert Maershal Error",err.Error()))
		}

		err = st.PutState(compositeKey,bytestr)
		if err != nil {
			return s.Success(resultMapMsg(false,21210,"Can not put state",err.Error()))
		}
	}

	response := resultMapMsg(true,9999,"Success func4 ","")
	st.SetEvent("func4",[]byte(args[0]))
	return s.Success(response)
}

func (c *Chaincode) func5(st s.ChaincodeStubInterface, args []string) p.Response {
	function, _ := st. GetFunctionAndParameters()
	log.Print("====================== " + function + " ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21300,"Incorrect number of argumnets 1.",""))
	}

	log.Print("Request Argument : " + args[0])
	var data map[string][]Struct2
	err := json.Unmarshal([]byte(args[0]), &data)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21301,"Filed to decode Json of: "+ string(args[0]),"",))
	}
	Data2, exists := data["data"]
	if !exists {
		return s.Success(resultMapMsg(false,21012,"Is not refund Data",""))
	}
	var eventData []Struct2
	for i:=0; i < len(Data2); i++ {
		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMapMsg(false,21307,"Can Not Create CompositeKey",err.Error()))
		}
		oldByteData, err := st.GetState(compositeKey)
		if err != nil {
			return s.Success(resultMapMsg(false,21308,"Get State Error",err.Error()))
		} else if oldByteData == nil {
			return s.Success(resultMapMsg(false,21309,"Not Refund Data",err.Error()))
		}

		var oldData Struct2
		err = json.Unmarshal(oldByteData,&oldData)
		if err != nil {
			return s.Success(resultMapMsg(false,21310,"Convert Maershal Error",err.Error()))
		}
		// 새로운 환불 데이터와 기존 환불 데이터 병합
		mergeMetaDate(&Data2[i], &oldData)

		// 병합한 환불 데이터 등록
		bytestr, err := json.Marshal(oldData)
		if err != nil {
			return s.Success(resultMapMsg(false,21311,"Convert Maershal Error",err.Error()))
		}
		
		log.Print("Change Data : " + string(bytestr))
		err = st.PutState(compositeKey,bytestr)
		if err != nil {
			return s.Success(resultMapMsg(false,21312,"Can not put state",err.Error()))
		}
		if function == "5" {
			eventData = append(eventData, oldData)	
		}
	}
	
	response := resultMapMsg(true,9999,"Success "+function,"")
	if function == "5" {
		bytestr, err := json.Marshal(eventData)
		if err != nil {
			return s.Success(resultMapMsg(false,21313,"Convert Maershal Error",err.Error()))
		}
		st.SetEvent("5Event",bytestr)
	}
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func6(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func6 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMapMsg(false,21400,"Incorrect number of argumnets 1.",[]string{}))
	}
		
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMapMsg(false,21401,"Filed to decode Json of: "+ string(args[0]),""))
	}

	flag := true
	if len(search.Rfnd_aplct_dates) != 0 {
		flag = false
	}

	var b bytes.Buffer
	var queryIterator s.StateQueryIteratorInterface
	// Rich Query 실행
	if flag == false {
		queryIterator, err := st.GetStateByPartialCompositeKey()
		if err != nil {
			return s.Success(resultMapMsg(false,21402,"query operation on private data failed. Error accessing state",[]string{}))
		}
		defer queryIterator.Close()
	} else if flag == true {
		queryString := fmt.Sprintf("{\"selector\":{\"1\":\"%s\"}}", search.1)
		b.WriteString(queryString)

		queryIterator, err = st.GetQueryResult(b.String())
		if err != nil {
			return s.Success(resultMapMsg(false,21403,"query operation on private data failed. Error accessing state",[]string{}))
		}
		defer queryIterator.Close()
	}
	
	var struct2InfoList []Struct2

	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMapMsg(false,21404,"query operation on private data failed. Error accessing state",[]string{}))
		}

		var struct2Info Struct2
		err := json.Unmarshal(response.Value, &struct2Info)
		if err != nil {
			return s.Success(resultMapMsg(false,21405,"Unmarshal Error %s, Convert Marshal",[]string{}))
		}
		struct2InfoList = append(struct2InfoList, struct2Info) //responce : key-value
	}
	
	response := resultMapMsg(true,9999,"Success func6",struct2InfoList)
	return s.Success(response)
}


func mergeMetaDate(src interface{}, dst interface{}) {
	// src 의 구조체 선언
	log.Printf("dat=", dst)
	dstTarget := reflect.ValueOf(dst)
	dstElements := dstTarget.Elem()

	// src 의 MetaData 가져오기
	srcTarget := reflect.ValueOf(src)
	srcElements := srcTarget.Elem()
	
	log.Printf("beforer=", dst)
	// src 의 MetaData nil 체크
	for i := 0; i < srcElements.NumField(); i++ {
		// src 의 MetaData nil 아니면 src 에 저장
		if srcElements.Field(i).Interface() != "" {
			dstElements.Field(i).Set(reflect.Value(srcElements.Field(i)))
		}
	}
	log.Printf("after=", dst)

}
/* 결과 정보 만드느 함수 */
func resultMapMsg(b bool, code int, msg string, data interface{}) []byte {
	var respMsg RespMsg
	respMsg.Success = b
	respMsg.Code	= code
	respMsg.Msg		= msg

	resultMap :=  make(map[string]interface{})
	resultMap["result"] = respMsg
	resultMap["data"] = data
	byteData, _ := json.Marshal(resultMap)
	
	return byteData
}
// main
func main() {
	err := s.Start(new(Chaincode))
	if err != nil {
		log.Println("Error create new Chaincode: ", err)
	}
}