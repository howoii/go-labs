// test how the tag `omitempty` affects json marshaller
package main

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Data struct {
	Str   string  `json:"str"`
	Id    int64   `json:"id"`
	SubId int     `json:"sub_id"`
	Value float32 `json:"value"`
}

func (e Data) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("str", e.Str)
	enc.AddInt64("id", e.Id)
	enc.AddInt("sub_id", e.SubId)
	enc.AddFloat32("value", e.Value)
	return nil
}

type DataDup struct {
	Str   string  `json:"str,omitempty"`
	Id    int     `json:"id,omitempty"`
	SubId int     `json:"sub_id,omitempty"`
	Value float32 `json:"value,omitempty"`
}

func (e DataDup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("str", e.Str)
	enc.AddInt("id", e.Id)
	enc.AddInt("sub_id", e.SubId)
	enc.AddFloat32("value", e.Value)
	return nil
}

func main() {
	logger := zap.NewExample()
	d1 := Data{
		Str:   "d1",
		Id:    35876901400,
		Value: 1.1,
	}
	d2 := DataDup{
		Str:   "d2",
		Id:    0,
		Value: 2.2,
	}
	fmt.Printf("d1: %v\n", d1)
	fmt.Printf("d2: %v\n", d2)
	s1, _ := json.Marshal(d1)
	s2, _ := json.Marshal(d2)
	fmt.Printf("jsonMarshal d1: %s\n", string(s1))
	fmt.Printf("jsonMarshal d2: %s\n", string(s2))
	logger.Debug("notOmit", zap.Object("data", d1))
	logger.Debug("omitted", zap.Object("data", d2))
}
