package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

type Namer interface {
	GetName() string
}

type accountMatchLikeInfo struct {
	MatchID  uint64 `db:"match_id" json:"match_id"`
	LikedCnt uint32 `db:"liked_cnt" json:"liked_cnt"`
}

type accountInfo struct {
	AuthorID         uint64               `db:"author_id" json:"author_id"`
	WorkshopID       uint32               `db:"workshop_id" json:"workshop_id"`
	WorkshopName     string               `db:"workshop_name" json:"workshop_name"`
	PerformanceCost  uint64               `db:"performance_cost" json:"performance_cost"`
	WorkshopSettings []byte               `db:"workshop_settings" json:"workshop_settings"`
	UpdateTime       int64                `db:"update_time" json:"update_time"`
	LikeInfo         accountMatchLikeInfo `db:"like_info" json:"like_info"`
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

type typeFiled struct {
	Name  string
	Index []int
	Typ   reflect.Type
}

type byIndex []typeFiled

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].Index {
		if k >= len(x[j].Index) {
			return false
		}
		if xik != x[j].Index[k] {
			return xik < x[j].Index[k]
		}
	}
	return len(x[i].Index) < len(x[j].Index)
}

func TestReflectStruct(t *testing.T) {
	info := accountInfo{}
	tp := reflect.TypeOf(info)
	var fields []typeFiled

	var curr, next []typeFiled
	next = append(next, typeFiled{Typ: tp})
	for len(next) > 0 {
		curr, next = next, curr[:0]
		for _, f := range curr {
			for i := 0; i < f.Typ.NumField(); i++ {
				index := make([]int, len(f.Index)+1)
				copy(index, f.Index)
				index[len(f.Index)] = i

				ff := f.Typ.Field(i)
				field := typeFiled{
					Name:  ff.Name,
					Index: index,
					Typ:   ff.Type,
				}
				fields = append(fields, field)

				if ff.Type.Kind() == reflect.Struct {
					next = append(next, field)
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].Name != x[j].Name {
			return x[i].Name < x[j].Name
		}
		if len(x[i].Index) != len(x[j].Index) {
			return len(x[i].Index) < len(x[j].Index)
		}
		return byIndex(x).Less(i, j)
	})

	for _, f := range fields {
		fmt.Println(f.Name, f.Index)
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

func typeSwitch(i interface{}) {
	switch v := i.(type) {
	case []int:
		fmt.Println(v, v == nil)
	case nil:
		fmt.Println("nil")
	default:
		fmt.Println("unexpected")
	}
}

func TestReflectNilType(t *testing.T) {
	var is []int
	typeSwitch(is)
	typeSwitch(nil)
	typeSwitch([]int{})
}
