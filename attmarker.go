package rootio

import (
	"bytes"
	"reflect"
)

type attmarker struct {
	Color int16
	Style int16
	Width float32
}

func (a *attmarker) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	}

	err = dec.readBin(a)
	if err != nil {
		return err
	}

	err = dec.checkByteCount(pos, bcnt, start, "TAttMarker")
	return err
}

//
func init() {
	f := func() reflect.Value {
		o := &attmarker{}
		return reflect.ValueOf(o)
	}
	Factory.db["TAttMarker"] = f
	Factory.db["*rootio.attmarker"] = f
}

// ifaces
var _ ROOTUnmarshaler = (*attmarker)(nil)
