package rootio

import (
	"bytes"
	"reflect"
)

type attaxis struct {
	Ndivs       int32   // number of division(10000*n3 + 100*n2 + n1)
	AxisColor   int16   // color of the line axis
	LabelColor  int16   // color of labels
	LabelFont   int16   // font for labels
	LabelOffset float32 // offset of labels
	LabelSize   float32 // size of labels
	TickLength  float32 // length of tick marks
	TitleOffset float32 // offset of axis title
	TitleSize   float32 // size of axis title
	TitleColor  int16   // color of axis title
	TitleFont   int16   // font for axis title
}

func (a *attaxis) UnmarshalROOT(data *bytes.Buffer) error {
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

	err = dec.checkByteCount(pos, bcnt, start, "TAttAxis")
	return err
}

//
func init() {
	f := func() reflect.Value {
		o := &attaxis{}
		return reflect.ValueOf(o)
	}
	Factory.db["TAttAxis"] = f
	Factory.db["*rootio.attaxis"] = f
}

// ifaces
var _ ROOTUnmarshaler = (*attaxis)(nil)
