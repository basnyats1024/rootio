package rootio

import (
	"bytes"
	"reflect"
)

// ObjArray is an array of Objects.
type objarray struct {
	name   string
	lbound int // lower bound of the array
	last   int // last element in array containing an object
	slice  []Object
}

func (arr *objarray) Class() string {
	return "TObjArray"
}

// Name returns the name of the instance
func (arr *objarray) Name() string {
	return arr.name
}

// Title returns the title of the instance
func (arr *objarray) Title() string {
	return arr.name
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *objarray) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	} else {
		myprintf("named: %v %v %v\n", vers, pos, bcnt)
	}

	if vers > 2 {
		var obj Object
		err = dec.readObject(&obj)
		if err != nil {
			return err
		}
	}
	if vers > 1 {
		err = dec.readString(&arr.name)
		if err != nil {
			return err
		}
	}

	nobjs := 0
	err = dec.readInt32(&nobjs)
	if err != nil {
		return err
	}
	arr.slice = make([]Object, 0, nobjs)

	err = dec.readInt32(&arr.lbound)
	if err != nil {
		return err
	}

	arr.last = -1

	for i := 0; i < nobjs; i++ {
		var obj Object
		err = dec.readObjectAny(&obj)
		if err != nil {
			return err
		}
		arr.slice = append(arr.slice, obj)
		arr.last = i
	}
	panic("not implemented")
	err = dec.checkByteCount(pos, bcnt, start, "TObjArray")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &objarray{
			slice: make([]Object, 0),
		}
		return reflect.ValueOf(o)
	}
	Factory.db["TObjArray"] = f
	Factory.db["*rootio.objarray"] = f
}

// check interfaces
var _ Object = (*objarray)(nil)
var _ ROOTUnmarshaler = (*objarray)(nil)

// EOF
