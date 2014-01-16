package rootio

import (
	"bytes"
	"reflect"
)

// tlist is a list of Objects.
type tlist struct {
	name  string // name of the collection
	slice []Object
}

func (lst *tlist) Class() string {
	return "TList"
}

// Name returns the name of the instance
func (lst *tlist) Name() string {
	return lst.name
}

// Title returns the title of the instance
func (lst *tlist) Title() string {
	return lst.name
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (lst *tlist) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	myprintf(">>> unmarshal tlist...\n")
	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	} else {
		myprintf("tlist: version=%v pos=%v count=%v\n", vers, pos, bcnt)
	}

	if vers > 3 {
		var obj Object
		err = dec.readObject(&obj)
		if err != nil {
			return err
		}

		err = dec.readString(&lst.name)
		if err != nil {
			return err
		}

		nobjs := 0
		err = dec.readInt32(&nobjs)
		if err != nil {
			return err
		}
		lst.slice = make([]Object, 0, nobjs)

		for i := 0; i < nobjs; i++ {
			var obj Object
			err = dec.readObject(&obj)
			if err != nil {
				return err
			}
			var nbig int32
			var nch uint8
			err = dec.readBin(&nch)
			if vers > 4 && nch == 255 {
				dec.readInt32(&nbig)
			} else {
				nbig = int32(nch)
			}
			bopt := make([]byte, int(nbig))
			err = dec.readBin(&bopt)
			if err != nil {
				return err
			}
			//readOption := string(bopt)
			if obj != nil {
				lst.slice = append(lst.slice, obj)
			}
		}
		err = dec.checkByteCount(pos, bcnt, start, "TList")
		return err
	}

	// process old versions when TList::Streamer was in TCollection::Streamer
	if vers > 2 {
		var obj Object
		err = dec.readObject(&obj)
		if err != nil {
			return err
		}
	}
	if vers > 1 {
		err = dec.readString(&lst.name)
		if err != nil {
			return err
		}
	}

	nobjs := 0
	err = dec.readInt32(&nobjs)
	if err != nil {
		return err
	}
	lst.slice = make([]Object, 0, nobjs)

	for i := 0; i < nobjs; i++ {
		var obj Object
		err = dec.readObject(&obj)
		if err != nil {
			return err
		}
		if obj != nil {
			lst.slice = append(lst.slice, obj)
		}
	}

	panic("not implemented")
	err = dec.checkByteCount(pos, bcnt, start, "TList")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &tlist{
			slice: make([]Object, 0),
		}
		return reflect.ValueOf(o)
	}
	Factory.db["TList"] = f
	Factory.db["*rootio.tlist"] = f
}

// check interfaces
var _ Object = (*tlist)(nil)
var _ ROOTUnmarshaler = (*tlist)(nil)

// EOF
