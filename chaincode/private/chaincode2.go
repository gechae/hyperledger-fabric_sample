package main

/*
	"fmt"
	"reflect"
	
	"strconv"
*/
import (
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"bytes"
	"reflect"
	"time"
	"strings"

	s "github.com/hyperledger/fabric/core/chaincode/shim"
	p "github.com/hyperledger/fabric/protos/peer"
)


type Chaincode struct {
}


type struct1 struct {

}

type struct2 struct {

}
type struct3 struct {

}

type struct4 struct {

}

type struct5 struct {

}

type struct6 struct {

}

type struct7 struct{

}

type struct8 struct {

// Search 조건 구조체
type SearchCond struct {
	dates_from		string			`json:",omitempty"`
	dates_to		string			`json:",omitempty"`

	Page_size			string			`json:"page_size"`
	Book_mark			string			`json:"book_mark"`
	
	amt			json.Number		`type:"integer" json:",omitempty"`
}
// 결과 정보 구조체
type RespMsg struct {
	Success		bool			`json:"success"`
	Code		int 			`json:"code"`
	Msg			string			`json:"msg"`
	Data		interface{}		`json:"data,omitempty"`
}
// pagination 결과 정보 구조체
type DataStrArry struct {
	InfoArry 			interface{}	`json:"Info"`
	Paginataion 			interface{}`json:"pagination"`
} 
func (c *Chaincode) Init(st s.ChaincodeStubInterface) p.Response {
	return s.Success(nil)
}

func (c *Chaincode) Invoke(st s.ChaincodeStubInterface) p.Response {
	function, args := st.GetFunctionAndParameters()
	
	switch function {
	case "1":
		return c.func1(st, args)
	case "2","3":
		 return c.func2(st, args, false)
	case "4":
		return c.func3(st, args, true)
	case "5","6","7","8","9":
		return c.func4(st, args)

	case "10":
		return c.func5(st, args)
	case "11":
		return c.func6(st, args)

	case "12":
		return c.func7(st, args)
	case "13":
		return c.func8(st, args)

	case "14":
		return c.func9(st, args)
	case "15":
		return c.func10(st, args)

	case "16":
		return c.func11(st, args)
	case "17":
		return c.func12(st, args)

	case "18":
		return c.func13(st, args)
	case "19":
		return c.func14(st, args)

	case "20":
		return c.func15(st, args)
	case "21":
		return c.func16(st, args)

	case "22":
		return c.func17(st, args)
	
	}

	response := resultMsg(false,8888,"Invalid Smart Contract function name.","")
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func1(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func1 ======================")
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21100,"Incorrect number of argumnets 1.",""))
	}

	log.Print("arg[0] :" +args[0])
	var struct1Arry []Struct1
	err := json.Unmarshal([]byte(args[0]), &struct1Arry)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,21101,"Filed to decode Json of: "+ string(args[0]),""))
	}


	for i:=0; i < len(struct1Arry); i++ {
		// vaildension key field check
		
		keyList:=strings.Split(token,"|")
		if len(keyList) != 4 {
			return s.Success(resultMsg(false,21107,"The keylist is not correct (4)",""))
		}
		// 복합키 생성
		log.Print("keyList 0 : " + keyList[0])
		log.Print("keyList 3 : " + keyList[3])
		compositeKey, err := st.CreateCompositeKey("dm",[]string{keyList[0],keyList[1],keyList[2],keyList[3]})
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,21108,"Can Not Create CompositeKey",err.Error()))
		}
		byteData, err := json.Marshal(struct1Arry[i])
		if err != nil {
			return s.Success(resultMsg(false,21109,"Convert Maershal Error",err.Error()))
		}

		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,21110,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func1","")
	log.Print(string(response))
	st.SetEvent("func1", response)
	return s.Success(response)
}

func (c *Chaincode) func2(st s.ChaincodeStubInterface, args []string, flag bool) p.Response {
	log.Print("====================== func2 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21201,"Incorrect number of argumnets 1.",""))
	}
	funct,_ := st.GetFunctionAndParameters()
	// 검색조건 데이터 Unmarshal
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,21202,"Can not Marshal",err.Error())))
		return s.Success(resultMsg(false,21202,"Can not Marshal",err.Error()))
	}

	// 검색 시작 일자 
	stDate, err := strconv.Atoi()
	if err != nil {
		return s.Success(resultMsg(false,21205,"Can not Convent Atoi",err.Error()))
	}
	// 검색 종료 일자
	enDate, err := strconv.Atoi()
	if err != nil {
		return s.Success(resultMsg(false,21206,"Can not Convent Atoi",err.Error()))
	}

	// pagination 에 필요한 데이터 
	var respSearch SearchCond
	var Struct1Arry []Struct1
	pageSize := 5000
	bookMark := ""
	if search.Page_size != "" {
		pageSize,_ = strconv.Atoi(search.Page_size)
	}
	if search.Book_mark != "" {
		bookMark = search.Book_mark
	}
	// pageSize Count 용
	_cnt := 0
	for i:=stDate; i <= enDate; {
		var b bytes.Buffer
		var queryRespMeta *p.QueryResponseMetadata
		var queryIterator s.StateQueryIteratorInterface
	
		// Create Query String, 검색 조건 query 생성
		idxNm:=1
		queryString := fmt.Sprintf("{\"selector\":{\"dates\":\"%s\"", strconv.Itoa(i))
		b.WriteString(queryString)
		switch funct {
		case "2","3": 
			if len(search.1) != 0 {
				queryString = fmt.Sprintf(",\"search2\":\"%s\"", )
				b.WriteString(queryString)
				idxNm=2
			}
			if len(search.2) != 0 {
				queryString = fmt.Sprintf(",\"search3\":\"%s\"", )
				b.WriteString(queryString)
				if idxNm != 1 {
					idxNm+=3
				}else {
					idxNm=3
				}
			}
		case "4":
			if len(search.amt.String()) !=0 {
				_ntpdAmt,_ := search.amt.Int64()
				queryString = fmt.Sprintf(",\"amt\":{\"$gt\":%d}", _ntpdAmt)
				b.WriteString(queryString)
			}else {
				queryString = fmt.Sprintf(",\"amt\":{\"$gt\":%d}", 0)
				b.WriteString(queryString)
			}
			idxNm=6
			if len(search.4) !=0 {
				queryString = fmt.Sprintf(",\"search4\":\"%s\"", search.4)
				b.WriteString(queryString)
					idxNm=7
			}
		}
		queryString = fmt.Sprintf("}, \"use_index\":[\"_design/%d\", \"%d\"]}",idxNm,idxNm)
		b.WriteString(queryString)
		log.Print("queryString : " + b.String())
		// Rich Query 실행 ( Pagination )
		queryIterator, queryRespMeta, err = st.GetQueryResultWithPagination(b.String(),int32(pageSize),bookMark)
		if err != nil {
			log.Print(string(resultMsg(false,21207,"Can not query",err.Error())))
			break
		}
		defer queryIterator.Close()

		logPrint := fmt.Sprintf("[LOG] Page Count : " , queryRespMeta.FetchedRecordsCount)
		log.Print(logPrint)
		log.Print("[LOG] BookMark : " + queryRespMeta.Bookmark )
		// 현재 검색된 Page Count
		var count = strconv.Itoa(int(queryRespMeta.FetchedRecordsCount))
		if count != "0" {
			log.Print("PageSize : " + respSearch.Page_size)
			respSearch.Page_size = count
		}
		// 현재 검색된 Book Mark
		if queryRespMeta.Bookmark != "nil" {
			log.Print("BookMark : " + queryRespMeta.Bookmark)
			respSearch.Book_mark = queryRespMeta.Bookmark
			log.Print("respSearch.Book_mark: ", respSearch.Book_mark)
		}
		
		for queryIterator.HasNext() {
			response, iterErr := queryIterator.Next()
			if iterErr != nil {
				return s.Success(returnData(21208,"No search query Iterateor",Struct1Arry,respSearch,false,flag))
			}

			var struct1 Struct1
			err := json.Unmarshal(response.Value, &struct1)
			if err != nil {
				return s.Success(returnData(21209,"Can not convert UnMarshaling",Struct1Arry,respSearch,false,flag))
			}

			Struct1Arry = append(Struct1Arry, struct1) //responce : key-value
			_cnt = _cnt+1
		}
		// 시작일과 종료일 비교
		stDateStr := strconv.Itoa(i)
		log.Print("stDateStr : " + stDateStr)
		fromTm := parseTime(strconv.Itoa(i))
		log.Print("From Time : " + fromTm.String())
		toTm := parseTime(search.dates_to)
		log.Print("To Time : " + toTm.String())
		if fromTm != toTm {
			t:=fromTm.AddDate(0,0,1).String()[:10]
			i, err = strconv.Atoi(strings.ReplaceAll(t,"-",""))
			if err != nil {
				return s.Success(returnData(21210,"Can not Convent Atoi", Struct1Arry,respSearch,false,flag))
			}
		}else {
			break
		}
		// pageSize가 검색 카운터와 동일하면 roof 종료
		if _cnt == pageSize {
			break
		}
	}
	
	response := returnData(9999,"Success func2",Struct1Arry,respSearch,false,flag)
	return s.Success(response)
}

// pagiantion 결괴 리턴 함수
func returnData(code int, msg string, args1 ,args2 interface{}, flag, b bool)  []byte{
	var response []byte
	if b == true {
		dataRelt := DataStrArry{args1,args2}
		response = resultMsg(flag,code,msg,dataRelt)
	} else {
		response = resultMsg(flag,code,msg,args1)
	}
	return response
}

func (c *Chaincode) func3(st s.ChaincodeStubInterface, args []string) p.Response{
	log.Print("====================== func3 ======================")
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21300,"Incorrect number of argumnets 1.",""))
	}
	
	log.Print("arg[0]" + args[0])
	// 새로 들어온 데이터 Unmarshal
	var struct1 []Struct1
	err := json.Unmarshal([]byte(args[0]), &struct1)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,21301,"Filed to decode Json of: "+ string(args[0]),""))
	}


	for i:=0; i < len(struct1); i++ {
		log.Print(info)
		
		keyList:=strings.Split(struct1[i].token,"|")
		if len(keyList) != 4 {
			return s.Success(resultMsg(false,21307,"The keylist is not correct (4)",""))
		}
		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,21308,"Can Not Create CompositeKey",err.Error()))
		}
		// 기존  데이터 가져오기, 없으면 등록
		oldByteData, err := st.GetState(compositeKey)
		if err != nil {
			return s.Success(resultMsg(false,21309,"Get State Error",err.Error()))
		} else if oldByteData == nil {
			bytestr, err := json.Marshal(struct1[i])
			if err != nil {
				return s.Success(resultMsg(false,21310,"Convert Maershal Error",err.Error()))
			}
			err = st.PutState(compositeKey, bytestr)
			if err != nil {
				return s.Success(resultMsg(false,21311,"Can not put state",err.Error()))
			}
			continue
		}

		var oldData PassInfoStruct
		err = json.Unmarshal(oldByteData,&oldData)
		if err != nil {
			return s.Success(resultMsg(false,21312,"Convert Maershal Error",err.Error()))
		}
		// 새로운  데이터와 기존  데이터 병합
		mergeMetaDate(&struct1[i], &oldData)

		// 병합한  데이터 등록
		bytestr, err := json.Marshal(oldData)
		if err != nil {
			return s.Success(resultMsg(false,21311,"Convert Maershal Error",err.Error()))
		}
		
		log.Print("Change Data : " + string(bytestr))
		err = st.PutState(compositeKey,bytestr)
		if err != nil {
			return s.Success(resultMsg(false,21312,"Can not put state",err.Error()))
		}
	}
	response := resultMsg(true,9999,"Success func3","")
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func4(st s.ChaincodeStubInterface, args []string) p.Response{
	log.Print("====================== func4 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21400,"Incorrect number of argumnets 1.",""))
	}

	log.Print("arg[0] :" +args[0])
	var struct2 []Struct2
	err := json.Unmarshal([]byte(args[0]), &struct2)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,21401,"Filed to decode Json of: "+ string(args[0]),""))
	}
	

	for i:=0; i < len(struct2); i++ {
		// vaildension key field check

		byteData, err := json.Marshal(struct2[i])
		if err != nil {
			return s.Success(resultMsg(false,21403,"Convert Maershal Error",err.Error()))
		}
		log.Print("SetPassInfo : " + string(byteData))
		err = st.PutState(struct2[i].1,byteData)
		if err != nil {
			return s.Success(resultMsg(false,21404,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func4","")
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func5(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func5 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21500,"Incorrect number of argumnets 1.",""))
	}
	// 검색조건 데이터 가져오기
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	var struct2 []Struct2
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,21501,"Can not Marshal",struct2)))
		return s.Success(resultMsg(false,21502,"Can not Marshal",struct2))
	}

	// 월정산 집게 데이터가져요기
	var struct2 Struct2

	err = json.Unmarshal(byteData, &struct2)
	if err != nil {
		return s.Success(resultMsg(false,9999,"Can not Marshal",struct2))
	}
	fadjSumMstArry = append(struct2, struct2) //responce : key-value 

	response := resultMsg(true,9999,"Success func5",struct2)
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func6(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func6 ======================")
	
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21600,"Incorrect number of argumnets 1.",""))
	}
	//  데이터 Unmarshal
	log.Print("arg[0] :" +args[0])
	var struct3 []Struct3
	err := json.Unmarshal([]byte(args[0]), &struct3)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,21601,"Filed to decode Json of: "+ string(args[0]),""))
	}

	for i:=0; i < len(struct3); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,21604,"Can Not Create CompositeKey",err.Error()))
		}

		byteData, err := json.Marshal(struct3[i])
		if err != nil {
			return s.Success(resultMsg(false,21605,"Convert Maershal Error",err.Error()))
		}
		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,21606,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func6","")
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func7(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func7 ======================")

	var struct3Arry []Struct3
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21700,"Incorrect number of argumnets 1.",""))
	}
	// 검색조건 데이터 unmarshal
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,21701,"Can not Marshal",struct3Arry)))
		return s.Success(resultMsg(false,21701,"Can not Marshal",struct3Arry))
	}
	
	// 검색조건 Key 값으로 데이터 가져오기
	queryIterator,err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMsg(false,21703,"Can not query",struct3Arry))
	}
	defer queryIterator.Close()
	
	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMsg(false,21704,"No search query Iterateor",struct3Arry))
		}
		var struct3 Struct3
		err := json.Unmarshal(response.Value, &struct3)
		if err != nil {
			return s.Success(resultMsg(false,21705,"Can not convert UnMarshaling",struct3Arry))
		}
		
		struct3Arry = append(struct3Arry, struct3) //responce : key-value
	}

	response := resultMsg(true,9999,"Success func7",struct3Arry)
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func8(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func8 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21800,"Incorrect number of argumnets 1.",""))
	}

	log.Print("arg[0] :" +args[0])
	var struct4rry []Struct4
	err := json.Unmarshal([]byte(args[0]), &struct4rry)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,21801,"Filed to decode Json of: "+ string(args[0]),""))
	}

	for i:=0; i < len(struct4rry); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,21804,"Can Not Create CompositeKey",err.Error()))
		}

		byteData, err := json.Marshal(struct4rry[i])
		if err != nil {
			return s.Success(resultMsg(false,21805,"Convert Maershal Error",err.Error()))
		}

		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,21806,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func8","")
	log.Print(string(response))
	return s.Success(response)
}

func (c *Chaincode) func9(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func9 ======================")

	var struct4Arry []Struct4

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,21900,"Incorrect number of argumnets 1.",""))
	}
	// 검색조건 데이터 가져오기
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,21901,"Can not Marshal",struct4Arry)))
		return s.Success(resultMsg(false,21901,"Can not Marshal",struct4Arry))
	}
	
	// 검색조건 Key 값으로  데이터 가져오기
	queryIterator,err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMsg(false,21902,"Can not query",struct4Arry))
	}
	defer queryIterator.Close()
	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMsg(false,21903,"No search query Iterateor",struct4Arry))
		}
		var struct4 Struct4
		err := json.Unmarshal(response.Value, &struct4)
		if err != nil {
			return s.Success(resultMsg(false,21904,"Can not convert UnMarshaling",struct4Arry))
		}
		
		struct4Arry = append(struct4Arry, struct4) //responce : key-value
	}

	response := resultMsg(true,9999,"Success func9",struct4Arry)
	log.Print(string(response))
	st.SetEvent("func9", response)
	return s.Success(response)
}

func (c *Chaincode) func10(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func10 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22000,"Incorrect number of argumnets 1.",""))
	}
	// 신대구/부산 고속도로 미납 수입금 데이터 Unmarshal
	log.Print("arg[0] :" +args[0])
	var struct5Arry []Struct5
	err := json.Unmarshal([]byte(args[0]), &struct5Arry)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,22001,"Filed to decode Json of: "+ string(args[0]),""))
	}


	for i:=0; i < len(struct5Arry); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,22004,"Can Not Create CompositeKey",err.Error()))
		}

		byteData, err := json.Marshal(struct5Arry[i])
		if err != nil {
			return s.Success(resultMsg(false,22005,"Convert Maershal Error",err.Error()))
		}
		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,22006,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func10","")
	log.Print(string(response))
	return s.Success(response)
}
/*
	신대구/부산 고속도로 미납 수입금 가져오기
*/
func (c *Chaincode) func11(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func11 ======================")

	var struct5Arry []Struct5

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22100,"Incorrect number of argumnets 1.",""))
	}
	// 검색조건 데이터 Unmarshal
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,22101,"Can not Marshal",struct5Arry)))
		return s.Success(resultMsg(false,22101,"Can not Marshal",struct5Arry))
	}

	// 검색조건 Key 값으로  데이터 가져오기
	queryIterator,err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMsg(false,22103,"Can not query",struct5Arry))
	}
	defer queryIterator.Close()

	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMsg(false,22104,"No search query Iterateor",struct5Arry))
		}
		var struct5 Struct5
		err := json.Unmarshal(response.Value, &struct5)
		if err != nil {
			return s.Success(resultMsg(false,22105,"Can not convert UnMarshaling",struct5Arry))
		}
		
		struct5Arry = append(struct5Arry, struct5) //responce : key-value
	}

	response := resultMsg(true,9999,"Success func11",struct5Arry)
	log.Print(string(response))

	return s.Success(response)
}

func (c *Chaincode) func12(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func12 ======================")

	
	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22200,"Incorrect number of argumnets 1.",""))
	}
	
	log.Print("arg[0] :" +args[0])
	var struct6Array []Struct6
	err := json.Unmarshal([]byte(args[0]), &struct6Array)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,22201,"Filed to decode Json of: "+ string(args[0]),""))
	}


	for i:=0; i < len(struct6Array); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,22204,"Can Not Create CompositeKey",err.Error()))
		}

		byteData, err := json.Marshal(struct6Array[i])
		if err != nil {
			return s.Success(resultMsg(false,22205,"Convert Maershal Error",err.Error()))
		}
		log.Print("func12 : " + string(byteData))
		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,22206,"Can not put state",err.Error()))
		}
	}


	response := resultMsg(true,9999,"Success func12","")
	log.Print(string(response))

	return s.Success(response)
}

func (c *Chaincode) func13(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func13 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22300,"Incorrect number of argumnets 1.",""))
	}

	var struct6Array []Struct6
	// 검색조건 데이터 Unmarshal
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,22301,"Can not Marshal",struct6Array)))
		return s.Success(resultMsg(false,22301,"Can not Marshal",struct6Array))
	}

	queryIterator,err := st.GetStateByPartialCompositeKey(})
	if err != nil {
		return s.Success(resultMsg(false,22303,"Can not query",struct6Array))
	}
	defer queryIterator.Close()

	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMsg(false,22304,"No search query Iterateor",struct6Array))
		}
		var struct6 Struct6
		err := json.Unmarshal(response.Value, &struct6)
		if err != nil {
			return s.Success(resultMsg(false,22305,"Can not convert UnMarshaling",struct6Array))
		}
		
		struct6Array = append(struct6Array, struct6) //responce : key-value
	}

	response := resultMsg(true,9999,"Success func13",struct6Array)
	log.Print(string(response))

	return s.Success(response)
}
/
func (c *Chaincode) func14(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func14 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22400,"Incorrect number of argumnets 1.",""))
	}
	// 연계할인 통행료  Unmarshal
	log.Print("arg[0] :" +args[0])
	var struct7Array []Struct7
	err := json.Unmarshal([]byte(args[0]), &struct7Array)
	if err != nil {
		log.Print("Filed to decode Json of :" + err.Error())
		return s.Success(resultMsg(false,22401,"Filed to decode Json of: "+ string(args[0]),""))
	}


	for i:=0; i < len(struct7Array); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,22404,"Can Not Create CompositeKey",err.Error()))
		}
		// 연계할인 통행료 저장
		byteData, err := json.Marshal(struct7Array[i])
		if err != nil {
			return s.Success(resultMsg(false,22405,"Convert Maershal Error",err.Error()))
		}

		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,22406,"Can not put state",err.Error()))
		}
	}

	response := resultMsg(true,9999,"Success func14","")
	log.Print(string(response))

	return s.Success(response)
}

func (c *Chaincode) func15(st s.ChaincodeStubInterface, args []string) p.Response {
	log.Print("====================== func15 ======================")

	var struct7Array []Struct7

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22500,"Incorrect number of argumnets 1.",""))
	}
	// 검색조건 데이터 Unmarshal
	log.Print("Request Argument : " + args[0])
	var search SearchCond
	err := json.Unmarshal([]byte(args[0]), &search)
	if err != nil {
		log.Print(string(resultMsg(false,22501,"Can not Marshal",struct7Array)))
		return s.Success(resultMsg(false,22501,"Can not Marshal",struct7Array))
	}


	// 검색조건 Key 값으로가져오기
	queryIterator,err := st.GetStateByPartialCompositeKey()
	if err != nil {
		return s.Success(resultMsg(false,22503,"Can not query",struct7Array))
	}
	defer queryIterator.Close()

	for queryIterator.HasNext() {
		response, iterErr := queryIterator.Next()
		if iterErr != nil {
			return s.Success(resultMsg(false,22504,"No search query Iterateor",struct7Array))
		}
		var struct7 Struct7
		err := json.Unmarshal(response.Value, &struct7)
		if err != nil {
			return s.Success(resultMsg(false,22505,"Can not convert UnMarshaling",struct7Array))
		}

		struct7Array = append(struct7Array, struct7) //responce : key-value
	}

	response := resultMsg(true,9999,"Success func15",struct7Array)
	log.Print(string(response))

	return s.Success(response)
}

func (c *Chaincode) func16(st s.ChaincodeStubInterface, args []string) p.Response{
	log.Print("====================== func16 ======================")

	if len(args) != 1 {
		log.Print("Incorrect number of argumnets 1.")
		return s.Success(resultMsg(false,22600,"Incorrect number of argumnets 1.",""))
	}
	// 비대면 데이터 Unmarshal
	log.Print("Request Argument : " + args[0])
	var struct8Arry []Struct8
	err := json.Unmarshal([]byte(args[0]), &struct8Arry)
	if err != nil {
		log.Print(string(resultMsg(false,22601,"Can not Marshal",[]string{})))
		return s.Success(resultMsg(false,22601,"Can not Marshal",[]string{}))
	}

	// vaildension key field check
	for i :=0; i < len(struct8Arry); i++ {

		// 복합키 생성
		compositeKey, err := st.CreateCompositeKey()
		if err != nil {
			log.Print("Can Not Create CompositeKey : " + err.Error())
			return s.Success(resultMsg(false,22603,"Can Not Create CompositeKey",[]string{}))
		}
		// 비대면 데이터 있는지 확인
		queryIterator,err := st.GetState(compositeKey)
		if err != nil {
			return s.Success(resultMsg(false,22604,"Can not query",[]string{}))
		}else if queryIterator != nil {
			return s.Success(resultMsg(false,9999,"Already Is Data exist",[]string{}))
		}
		// 비대면 정보 저장
		byteData, err := json.Marshal(struct8Arry[i])
		if err != nil {
			return s.Success(resultMsg(false,22605,"Convert Maershal Error",[]string{}))
		}

		err = st.PutState(compositeKey,byteData)
		if err != nil {
			return s.Success(resultMsg(false,22606,"Can not put state",[]string{}))
		}
	}
	response := resultMsg(true,9999,"Success func16",[]string{})
	log.Print(string(response))
	
	return s.Success(response)
}
/*
	전달 받은 시간 데이터 parse
*/
func parseTime(date string) time.Time {
	log.Print("get time string : " + date)
	test := time.Now().Format("20060102")
	log.Print("test : " + test)
	year,_:= strconv.Atoi(date[:4])
	month,_ := strconv.Atoi(date[4:6])
	day,_ := strconv.Atoi(date[6:])
	t := time.Date(year,time.Month(month),day,0,0,0,0, time.UTC)
	return t
}
/*
	새로운 데이터 와 이전 데이터 Merge 함수
*/
func mergeMetaDate(src interface{}, dst interface{}) {
	// src 의 구조체 선언
	log.Printf("dat=", dst)
	dstTarget := reflect.ValueOf(dst)
	dstElements := dstTarget.Elem()

	// src 의 MetaData 가져오기
	srcTarget := reflect.ValueOf(src)
	srcElements := srcTarget.Elem()

	log.Printf("before=", dst)
	// src 의 MetaData nil 체크
	for i := 0; i < srcElements.NumField(); i++ {
		// src 의 MetaData nil 아니면 src 에 저장
		if srcElements.Field(i).Interface() != "" {
			dstElements.Field(i).Set(reflect.Value(srcElements.Field(i)))
		}
	}
	log.Printf("after=", dst)

}
/*
	결과 데이터 만드는 함수
*/
func resultMsg(b bool, code int, msg string, data interface{}) []byte {
	
	respMsg := RespMsg{b, code, msg, data}

	byteData, err := json.Marshal(respMsg)
	if err != nil {
		log.Print("Convert Maershal Error : " + err.Error())
	}
	return byteData
}
// main
func main() {
	err := s.Start(new(Chaincode))
	if err != nil {
		log.Println("Error create new Chaincode: ", err)
	}
}
