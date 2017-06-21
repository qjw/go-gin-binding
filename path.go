// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"strings"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
)

type pathBinding struct{}

func (pathBinding) Name() string {
	return "path"
}

func (obj pathBinding) Bind(c * gin.Context, data interface{}) error {
	// return obj.read(c,data)
	value := reflect.ValueOf(data)
	if err := obj.read(c,value);err != nil{
		return err
	}
	return Validate(data)
}

func (obj pathBinding) read(c * gin.Context,val reflect.Value) error{
	// t := reflect.ValueOf(data).Type()
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Struct:
		//typ := t.Elem()
		//val := v.Elem()
		count := typ.NumField()
		for i := 0; i < count; i++ {
			typeField := typ.Field(i)
			structField := val.Field(i)

			// 只能是基本类型
			fmt.Println(typeField.Type.Kind())
			if _, ok := kindMapping[typeField.Type.Kind()]; !ok {
				return errors.New("path object invalid field type")
			}

			if !structField.CanSet() {
				continue
			}

			tag := typeField.Tag.Get("json")
			name := parseTag(tag)
			if name == "" {
				name = typeField.Name
			}
			if name == "-" {
				continue
			}

			value := c.Param(name)
			if err := setWithProperType(typeField.Type, value, structField); err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return obj.read(c,val.Elem())
	default:
		return errors.New("path object invalid type")
	}
}

var kindMapping = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
}

func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

//func (obj pathBinding) readFromStruct(t reflect.Type,c * gin.Context) {
//	count := t.NumField()
//	if count < 1{
//		log.Fatal("path object invalid field count")
//		return
//	}
//	for i := 0; i < count; i++ {
//		field := t.Field(i)
//		field_var := t.Elem().Field(i)
//		field_type := field.Type
//		if !field_var.CanSet() {
//			continue
//		}
//
//		// 只能是基本类型
//		if _, ok := kindMapping[field_type.Kind()]; ok {
//			log.Fatal("path object invalid field type")
//			return
//		}
//
//		tag := field.Tag.Get("json")
//		name := parseTag(tag)
//		if name == "" {
//			name = field.Name
//		}
//		if name == "-" {
//			continue
//		}
//
//		value := c.Param(name)
//		setWithProperType(field_type.Kind(),value,field)
//	}
//}