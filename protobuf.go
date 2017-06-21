// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/golang/protobuf/proto"

	"io/ioutil"
	"github.com/gin-gonic/gin"
)

type protobufBinding struct{}

func (protobufBinding) Name() string {
	return "protobuf"
}

func (protobufBinding) Bind(c * gin.Context, obj interface{}) error {

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	if err = proto.Unmarshal(buf, obj.(proto.Message)); err != nil {
		return err
	}

	//Here it's same to return validate(obj), but util now we cann't add `binding:""` to the struct
	//which automatically generate by gen-proto
	return nil
	//return validate(obj)
}
