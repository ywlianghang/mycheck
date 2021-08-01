package Stream

import (
	"fmt"
	"strconv"
	"time"
)

type StreamStruct struct{
	argss []interface{}
	Data []map[string]interface{}
}

//map函数
func (s StreamStruct) Filter(f func(interface{}) bool) StreamStruct {
	res := make([]map[string]interface{}, 0, len(s.Data))
	for _, item := range s.Data {
		if f(item){
			res = append(res, item)
		}
	}
	return StreamStruct{Data:res}
}
//fliter函数
//func (s StreamStruct) Filter(f func(int) bool) StreamStruct {
//	res := []int{}
//	for _, item := range s.data {
//		if f(item) {
//			res = append(res, item)
//		}
//	}
//	return StreamStruct{data:res}
//}

//func main() {
//	stream := Stream{[]int{1, 3, 5, 7, 8}}
//	fmt.Println(stream.data)
//	stream = stream.
//		Map(func(i int) int { return i + 1 }).
//		Filter(func(i int) bool { return i%2 == 0 }).
//		Filter(func(i int) bool { return i >5})
//	fmt.Println(stream.data)
//
//	>>>[1 3 5 7 8]
//[6 8]
//}
func (stre *StreamStruct) GetTimeDayArr (startTime,endTime string) (int64,error){
	// 转成时间戳
	timeLayout  := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startUnix,err := time.ParseInLocation(timeLayout,startTime,loc)
	if err != nil{
		return 0,err
	}
	endUnix,err := time.ParseInLocation(timeLayout,endTime,loc)
	if err != nil{
		return 0,err
	}
	//求相差天数
	dateDay :=	(endUnix.Unix() - startUnix.Unix()) / 86400
	return dateDay,nil
}

func (stre *StreamStruct) baseDateDispose(args interface{}) ([]string,error){
	var argsSlict []string
	//if len(args) >2{
	//	err := "If there are more than two input parameters, the calculation cannot be performed"
	//	return argsSlict,errors.New(err)
	//}
	//strArray := make([]string, len(args))
	//for i, arg := range args { strArray[i] = arg.(string) }
	//aa := strings.Join(strArray, "_")
	//c,err := stre.BaseDateTypeInterfaceToIntSlice(args)
	//fmt.Println(c,err)

	//for _,arg := range args{
	//	switch v := arg.(type){
	//	case string:
	//		fmt.Println(v)
	//	case int:
	//		fmt.Println(v)
	//	default:
	//		fmt.Println("params type not supported")
	//	}
	//}
	//	argsSlict = append(argsSlict,arg.(string))
	//	b := strings.Split(fmt.Sprintf("%v",v)," ")
	//	a,err := strconv.Atoi(b[1][:len(b[1])-1])
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	argsSlict = append(argsSlict,a)
	//}
	return argsSlict,nil
}
func (stre *StreamStruct) Add(args ...interface{}) (int,error){
	var intArry = make([]int,len(args))
	for i,v:= range args{
		d,err := strconv.Atoi(fmt.Sprintf("%v",v))
		if err != nil{return 0,err}
		intArry[i] = d
	}
	c := intArry[0] + intArry[1]
	return c,nil
}

//返回值计算的百分比%  a *100 / b
func (stre *StreamStruct) Percentage(args ...interface{}) (int,error){
	var intArry = make([]int,len(args))
	for i,v:= range args{
		d,err := strconv.Atoi(fmt.Sprintf("%v",v))
		if err != nil{return 0,err}
		intArry[i] = d
	}
	c := intArry[0] * 100 /intArry[1]
	return c,nil
}
