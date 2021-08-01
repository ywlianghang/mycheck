package Stream

import (
"errors"
"reflect"
"sort"
)

type Stream struct {
	ops  []op
	data []interface{}
}

type op struct {
	typ string
	fun reflect.Value
}

type sortbyfun struct {
	data []interface{}
	fun  reflect.Value
}

func (s sortbyfun) Len() int            { return len(s.data) }
func (s *sortbyfun) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s *sortbyfun) Less(i, j int) bool { return call(s.fun, s.data[i], s.data[j])[0].Bool() }

// New create a stream from a slice
func New(arr interface{}) (*Stream, error) {
	ops := make([]op, 0)
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Kind() == reflect.Slice || arrValue.Kind() == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			dataValue.Set(reflect.Append(dataValue, arrValue.Index(i)))
		}
	} else {
		return nil, errors.New("the type of arr parameter must be Array or Slice")
	}
	return &Stream{ops: ops, data: data}, nil
}

// Of create a stream from some values
func Of(args ...interface{}) (*Stream, error) {
	return New(args)
}

// Ints create a stream from some ints.
func Ints(args ...int64) (*Stream, error) {
	return New(args)
}

// Floats create a stream from some floats.
func Floats(args ...float64) (*Stream, error) {
	return New(args)
}

// Strings create a stream from some strings.
func Strings(args ...string) (*Stream, error) {
	return New(args)
}

// It create a stream from a iterator.itFunc: func(prev T) (next T,more bool)
func It(initValue interface{}, itFunc interface{}) (*Stream, error) {
	funcValue := reflect.ValueOf(itFunc)
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	prev := reflect.ValueOf(initValue)
	for {
		out := funcValue.Call([]reflect.Value{prev})
		dataValue.Set(reflect.Append(dataValue, out[0]))
		if !out[1].Bool() {
			break
		}
		prev = out[0]
	}
	return New(data)
}

// Gen create a stream by invoke genFunc. genFunc: func() (next T,more bool)
func Gen(genFunc interface{}) (*Stream, error) {
	funcValue := reflect.ValueOf(genFunc)
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	for {
		out := call(funcValue)
		dataValue.Set(reflect.Append(dataValue, out[0]))
		if !out[1].Bool() {
			break
		}
	}
	return New(data)
}

func (s *Stream) Reset() *Stream {
	s.ops = make([]op, 0)
	return s
}

//  Filter operation. filterFunc: func(o T) bool
func (s *Stream) Filter(filterFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(filterFunc)
	s.ops = append(s.ops, op{typ: "filter", fun: funcValue})
	return s
}

//  Map operation. Map one to one
// mapFunc: func(o T1) T2
func (s *Stream) Map(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "map", fun: funcValue})
	return s
}

// FlatMap operation. Map one to many
// mapFunc: func(o T1) []T2
func (s *Stream) FlatMap(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "flatMap", fun: funcValue})
	return s
}

// Sort operation. lessFunc: func(o1,o2 T) bool
func (s *Stream) Sort(lessFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(lessFunc)
	s.ops = append(s.ops, op{typ: "sort", fun: funcValue})
	return s
}

// Distinct operation. equalFunc: func(o1,o2 T) bool
func (s *Stream) Distinct(equalFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(equalFunc)
	s.ops = append(s.ops, op{typ: "distinct", fun: funcValue})
	return s
}

// Peek operation. peekFunc: func(o T)
func (s *Stream) Peek(peekFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(peekFunc)
	s.ops = append(s.ops, op{typ: "peek", fun: funcValue})
	return s
}

// Call operation. Call function with the data.
// callFunc: func()
func (s *Stream) Call(callFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(callFunc)
	s.ops = append(s.ops, op{typ: "call", fun: funcValue})
	return s
}

// Check operation. Check if should be continue process data.
// checkFunc: func(o []T) bool ,checkFunc must return if should be continue process data.
func (s *Stream) Check(checkFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(checkFunc)
	s.ops = append(s.ops, op{typ: "check", fun: funcValue})
	return s
}

// Limit operation.
func (s *Stream) Limit(num int) *Stream {
	if num < 0 {
		num = 0
	}
	funcValue := reflect.ValueOf(func() int { return num })
	s.ops = append(s.ops, op{typ: "limit", fun: funcValue})
	return s
}

// Skip operation.
func (s *Stream) Skip(num int) *Stream {
	if num < 0 {
		num = 0
	}
	funcValue := reflect.ValueOf(func() int { return num })
	s.ops = append(s.ops, op{typ: "skip", fun: funcValue})
	return s
}

// collect operation.
func (s *Stream) collect() []interface{} {
	result := s.data
	for _, op := range s.ops {
		if len(result) == 0 {
			break
		}
		switch op.typ {
		case "filter":
			temp := make([]interface{}, 0)
			each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
				if out[0].Bool() {
					temp = append(temp, it)
				}
				return true
			})
			result = temp
		case "peek":
			each(result, op.fun, emptyeachfunc)
		case "map":
			temp := make([]interface{}, 0)
			tempVlaue := reflect.ValueOf(&temp).Elem()
			each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
				tempVlaue.Set(reflect.Append(tempVlaue, out[0]))
				return true
			})
			result = temp
		case "flatMap":
			temp := make([]interface{}, 0)
			tempVlaue := reflect.ValueOf(&temp).Elem()
			each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
				for i := 0; i < out[0].Len(); i++ {
					tempVlaue.Set(reflect.Append(tempVlaue, out[0].Index(i)))
				}
				return true
			})
			result = temp
		case "aggMap":

		case "sort":
			sort.Sort(&sortbyfun{data: result, fun: op.fun})
		case "distinct":
			temp := make([]interface{}, 0)
			temp = append(temp, result[0])
			for _, it := range result {
				found := false
				for _, it2 := range temp {
					out := call(op.fun, it, it2)
					if out[0].Bool() {
						found = true
					}
				}
				if !found {
					temp = append(temp, it)
				}
			}
			result = temp
		case "limit":
			limit := int(call(op.fun)[0].Int())
			if limit > len(result) {
				limit = len(result)
			}
			temp := result
			result = temp[:limit]
		case "skip":
			skip := int(call(op.fun)[0].Int())
			if skip > len(result) {
				skip = len(result)
			}
			temp := result
			result = temp[skip:]
		case "call":
			call(op.fun)
		case "check":
			out := call(op.fun, result)
			if !out[0].Bool() {
				break
			}
		}
	}
	return result
}

// Exec operation.
func (s *Stream) Exec() {
	s.collect()
}

// ToSlice operation. targetSlice must be a pointer.
func (s *Stream) ToSlice(targetSlice interface{}) error {
	data := s.collect()
	targetValue := reflect.ValueOf(targetSlice)
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("target slice must be a pointer")
	}
	sliceValue := reflect.Indirect(targetValue)
	for _, it := range data {
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(it)))
	}
	return nil
}

// ForEach executes a provided function once for each array element,and terminate the stream.
// actFunc: func(o T)
func (s *Stream) ForEach(actFunc interface{}) {
	data := s.collect()
	each(data, reflect.ValueOf(actFunc), emptyeachfunc)
}

// AllMatch operation.
// matchFunc: func(o T) bool
func (s *Stream) AllMatch(matchFunc interface{}) bool {
	data := s.collect()
	allMatch := true
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if !out[0].Bool() {
			allMatch = false
			return false
		}
		return true
	})
	return allMatch
}

// AnyMatch operation. matchFunc: func(o T) bool
func (s *Stream) AnyMatch(matchFunc interface{}) bool {
	data := s.collect()
	anyMatch := false
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if out[0].Bool() {
			anyMatch = true
			return false
		}
		return true
	})
	return anyMatch
}

// NoneMatch operation. matchFunc: func(o T) bool
func (s *Stream) NoneMatch(matchFunc interface{}) bool {
	data := s.collect()
	noneMatch := true
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if out[0].Bool() {
			noneMatch = false
			return false
		}
		return true
	})
	return noneMatch
}

// Count operation.Return the count of elements in stream.
func (s *Stream) Count() int {
	return len(s.collect())
}

// Group operation. Group values by key.
// Premeter groupFunc: func(o T1) (key T2,value T3). Return map[T2]T3
func (s *Stream) Group(groupFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(groupFunc)
	result := make(map[interface{}][]interface{})
	rValue := reflect.ValueOf(result)

	for _, it := range data {
		out := call(funcValue, it)
		sliceValue := rValue.MapIndex(out[0])

		if !sliceValue.IsValid() || sliceValue.IsNil()  {
			sliceValue = reflect.ValueOf(make([]interface{}, 0))
		}
		value := reflect.Append(sliceValue, out[1])
		rValue.SetMapIndex(out[0], value)
	}
	return rValue.Interface()
}

// Max operation.lessFunc: func(o1,o2 T) bool
func (s *Stream) Max(lessFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(lessFunc)
	data := s.collect()
	var max interface{}
	if len(data) > 0 {
		max = data[0]
		for i := 1; i < len(data); i++ {
			out := call(funcValue, max, data[i])
			if out[0].Bool() {
				max = data[i]
			}
		}
	}
	return max
}

// Min operation.lessFunc: func(o1,o2 T) bool
func (s *Stream) Min(lessFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(lessFunc)
	data := s.collect()
	var min interface{}
	if len(data) > 0 {
		min = data[0]
		for i := 1; i < len(data); i++ {
			out := call(funcValue, data[i], min)
			if out[0].Bool() {
				min = data[i]
			}
		}
	}
	return min
}

// First operation. matchFunc: func(o T) bool
func (s *Stream) First(matchFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(matchFunc)
	for _, it := range data {
		out := call(funcValue, it)
		if out[0].Bool() {
			return it
		}
	}
	return nil
}

// Last operation. matchFunc: func(o T) bool
func (s *Stream) Last(matchFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(matchFunc)
	for i := len(data) - 1; i >= 0; i-- {
		it := data[i]
		out := call(funcValue, it)
		if out[0].Bool() {
			return it
		}
	}
	return nil
}

// Reduce operation. reduceFunc: func(r T2,o T) T2
func (s *Stream) Reduce(initValue interface{}, reduceFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(reduceFunc)
	result := initValue
	rValue := reflect.ValueOf(&result).Elem()
	for _, it := range data {
		out := call(funcValue, result, it)
		rValue.Set(out[0])
	}
	return result
}

//  eachfunc is the function for each method,return if should continue loop
type eachfunc func(int, interface{}, []reflect.Value) bool

//  emptyeachfunc the empty eachfunc, return true
var emptyeachfunc = func(int, interface{}, []reflect.Value) bool { return true }

func each(data []interface{}, fun reflect.Value, act eachfunc) {
	for i, it := range data {
		out := call(fun, it)
		if !act(i, it, out) {
			break
		}
	}
}

func call(fun reflect.Value, args ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(args))
	for i, a := range args {
		in[i] = reflect.ValueOf(a).Convert(fun.Type().In(i))
	}
	return fun.Call(in)
}

