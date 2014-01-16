package rootio

import (
	"bytes"
	"reflect"
)

// ObjArray is an array of Objects.
type objarray struct {
	named named
	slice []Object
}

func (arr *objarray) Class() string {
	return "TObjArray"
}

func (arr *objarray) Name() string {
	return arr.named.Name()
}

func (arr *objarray) Title() string {
	return arr.named.Title()
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *objarray) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
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
