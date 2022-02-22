package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	objectTemplate = `    %s:
      type: object
      description: %s`
	propertyTemplate = `        %s:
          type: %s
          example: %s`
	arrayRefTemplate = `        %s:
          type: array
          items:
            $ref: '#/components/schemas/%s'`
	arrayTemplate = `        %s:
          type: array
          items:
            type: %s`
	structTemplate = `        %s:
          $ref: '#/components/schemas/%s'`
)

type Property struct {
	Name    string
	Example string
	Typ     reflect.Type
}

type Object struct {
	Name       string
	Properties []Property
}

func getFieldProperty(f reflect.StructField) Property {
	p := Property{
		Typ: f.Type,
	}
	tagFields := strings.Split(f.Tag.Get("json"), ",")
	if len(tagFields) > 0 {
		p.Name = tagFields[0]
	} else {
		p.Name = f.Name
	}
	if example, ok := f.Tag.Lookup("example"); ok {
		p.Example = example
	}
	return p
}

func kindString(k reflect.Kind) string {
	switch k {
	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "string"
	default:
		return k.String()
	}
}

func formatProperty(p Property, w io.Writer) {
	k := p.Typ.Kind()
	switch k {
	case reflect.Struct:
		fmt.Fprintln(w, fmt.Sprintf(structTemplate, p.Name, p.Typ.Name()))
	case reflect.Slice:
		var et reflect.Type
		if p.Typ.Elem().Kind() == reflect.Ptr {
			et = p.Typ.Elem().Elem()
		} else {
			et = p.Typ.Elem()
		}
		if et.Kind() == reflect.Struct {
			fmt.Fprintln(w, fmt.Sprintf(arrayRefTemplate, p.Name, et.Name()))
		} else {
			fmt.Fprintln(w, fmt.Sprintf(arrayTemplate, p.Name, kindString(et.Kind())))
		}
	default:
		fmt.Fprintln(w, fmt.Sprintf(propertyTemplate, p.Name, kindString(k), p.Example))
	}
}

func formatObject(o *Object, w io.Writer) {
	fmt.Fprintln(w, fmt.Sprintf(objectTemplate, o.Name, o.Name))
	if len(o.Properties) == 0 {
		return
	}
	fmt.Fprintln(w, "      properties:")
	for _, p := range o.Properties {
		formatProperty(p, w)
	}
}

func main() {
	var objs []*Object
	visited := make(map[string]struct{})

	var ts []reflect.Type
	structs := []interface{}{
		//	add your structs
	}

	for _, v := range structs {
		t := reflect.TypeOf(v)
		visited[t.String()] = struct{}{}

		ts = append(ts, t)
		for len(ts) > 0 {
			t := ts[0]
			ts = ts[1:]

			o := &Object{
				Name: t.Name(),
			}
			objs = append(objs, o)
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				p := getFieldProperty(f)
				o.Properties = append(o.Properties, p)

				ft := f.Type
				if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
					ft = ft.Elem().Elem()
				} else if ft.Kind() == reflect.Slice {
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct {
					if _, ok := visited[ft.String()]; !ok {
						visited[ft.String()] = struct{}{}
						ts = append(ts, ft)
					}
				}
			}
		}
	}

	var w io.Writer
	f, err := os.OpenFile("./openapi.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("close file err: %v", err)
		}
	}()
	if err == nil {
		w = f
	} else {
		w = os.Stdout
		log.Printf("open file err: %v", err)
	}
	for _, o := range objs {
		formatObject(o, w)
	}
}
