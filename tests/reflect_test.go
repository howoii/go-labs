package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type Namer interface {
	GetName() string
}

var typeNamer = reflect.TypeOf((*Namer)(nil)).Elem()

func TestReflectType(t *testing.T) {
	var i interface{} = (*interface{})(nil)
	switch v := i.(type) {
	case *uint32:
		fmt.Println("*uint", v)
	case *interface{}:
		fmt.Println("interface", v)
	case *Namer:
		fmt.Println(reflect.TypeOf(v))
	default:
		t := reflect.TypeOf(i)
		fmt.Println(t.String())
	}

}

func TestReflectMethod(t *testing.T) {
	account := &AccountInfo{1, "ddd", "10086"}
	func(i interface{}) {
		t := reflect.TypeOf(i)
		v := reflect.ValueOf(i)
		fmt.Println(t.Kind())

		tt := reflect.TypeOf(t)
		fmt.Println(tt.Kind())

		if t.NumMethod() > 0 {
			m, ok := t.MethodByName("GetName")
			if ok {
				fmt.Println(m)
				vs := m.Func.Call([]reflect.Value{v})
				for _, v := range vs {
					fmt.Println(v.String())
				}
			}
		}
		if v.NumMethod() > 0 {
			m := v.Method(0)
			vs := m.Call(nil)
			for _, v := range vs {
				fmt.Println(v.Uint())
			}
		}
	}(account)
}

func TestReflectElem(t *testing.T) {
	sl := []int{}
	m := map[int]string{}
	f := func(i interface{}) {
		t := reflect.TypeOf(i)
		fmt.Println(t.Kind())
		e := t.Elem()
		fmt.Println(e.Kind())
	}
	f(sl)
	f(m)
}

func TestReflectNew(t *testing.T) {
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age"`,
		},
	})

	v := reflect.New(typ).Elem()
	v.Field(0).SetFloat(0.4)
	v.Field(1).SetInt(2)
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	fmt.Printf("value: %+v\n", s)
	fmt.Printf("json:  %s", w.Bytes())

	r := bytes.NewReader([]byte(`{"height":1.5,"age":10}`))
	if err := json.NewDecoder(r).Decode(s); err != nil {
		panic(err)
	}
	fmt.Printf("value: %+v\n", s)
}

func TestReflectInterface(t *testing.T) {
	fmt.Println(reflect.TypeOf(Namer(nil)))

	nilT := reflect.TypeOf((*Namer)(nil))
	fmt.Println(nilT.Kind())

	ele := nilT.Elem()
	fmt.Println(ele.Kind())

	acT := reflect.TypeOf((*AccountInfo)(nil))
	fmt.Println(acT.Kind())

	fmt.Println(acT.Implements(ele))
}

func TestReflectValue(t *testing.T) {
	a := AccountInfo{}
	v := reflect.ValueOf(a)
	fmt.Println(v.Kind())

	vs := reflect.ValueOf([]int{1, 2})
	fmt.Println(vs.Kind())
	is := vs.Index(0)
	fmt.Println(is.Int())
	as := is.Addr()
	fmt.Println(as.Kind())
	fmt.Println(as.Elem())
}
