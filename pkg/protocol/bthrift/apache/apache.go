/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package apache

import (
	"errors"
	"io"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/gopkg/protocol/thrift/apache"
)

func init() {
	// it makes github.com/cloudwego/gopkg/protocol/thrift/apache works
	apache.RegisterCheckTStruct(checkTStruct)
	apache.RegisterThriftRead(callThriftRead)
	apache.RegisterThriftWrite(callThriftWrite)
}

var errNotThriftTStruct = errors.New("not thrift.TStruct")

func checkTStruct(v interface{}) error {
	_, ok := v.(thrift.TStruct)
	if !ok {
		return errNotThriftTStruct
	}
	return nil
}

func callThriftRead(r io.ReadWriter, v interface{}) error {
	p, ok := v.(thrift.TStruct)
	if !ok {
		return errNotThriftTStruct
	}
	t, ok := r.(byteBuffer)
	if ok {
		in := NewBinaryProtocol(t)
		return p.Read(in)
	}
	in := thrift.NewTBinaryProtocol(apache.NewDefaultTransport(r), true, true)
	return p.Read(in)
}

func callThriftWrite(w io.ReadWriter, v interface{}) error {
	p, ok := v.(thrift.TStruct)
	if !ok {
		return errNotThriftTStruct
	}
	t, ok := w.(byteBuffer)
	if ok {
		out := NewBinaryProtocol(t)
		return p.Write(out)
	}
	out := thrift.NewTBinaryProtocol(apache.NewDefaultTransport(w), true, true)
	return p.Write(out)
}
