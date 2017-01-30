// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"reflect"
)

type tlist struct {
	name string
	objs []Object
}

func (li *tlist) Class() string {
	return "TList"
}

func (li *tlist) Name() string {
	if li.name == "" {
		return "TList"
	}
	return li.name
}

func (li *tlist) At(i int) Object {
	return li.objs[i]
}

func (li *tlist) Last() int {
	panic("not implemented")
}

func (li *tlist) Len() int {
	return len(li.objs)
}

func (li *tlist) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	myprintf("tlist-vers=%v\n", vers)

	if vers <= 3 {
		return fmt.Errorf("rootio: TList version too old (%d <= 3)", vers)
	}

	_ = r.ReadU32() // id
	bits := r.ReadU32()
	bits |= kIsOnHeap // by definition, de-serialized objects are on the heap
	if bits&kIsReferenced == 0 {
		_ = r.ReadU16()
	}
	li.name = r.ReadString()
	size := int(r.ReadI32())

	li.objs = make([]Object, size)

	for i := range li.objs {
		obj := r.ReadObjectAny()
		// obj := r.ReadObjectRef()
		if obj == nil {
			panic("nil obj ref") // FIXME(sbinet)
			return r.Err()
		}
		li.objs[i] = obj

		n := int(r.ReadU8())
		if n > 0 {
			opt := make([]byte, n)
			io.ReadFull(r.r, opt)
			// drop the option on the floor. // FIXME(sbinet)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TList")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tlist{}
		return reflect.ValueOf(o)
	}
	Factory.add("TList", f)
	Factory.add("*rootio.tlist", f)
}

var _ Object = (*tlist)(nil)
var _ Collection = (*tlist)(nil)
var _ SeqCollection = (*tlist)(nil)
var _ List = (*tlist)(nil)
var _ ROOTUnmarshaler = (*tlist)(nil)