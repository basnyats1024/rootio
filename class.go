package rootio

import (
	"bytes"
	"reflect"
)

type class struct {
	named   named
	members []Member // data members
	version int      // class version identifier
}

func (cls *class) Class() string {
	return "TClass"
}

// Name returns the name of the instance
func (cls *class) Name() string {
	return cls.named.Name()
}

// Title returns the title of the instance
func (cls *class) Title() string {
	return cls.named.Title()
}

// Members returns the list of members for this ROOT class
func (cls *class) Members() []Member {
	return cls.members
}

// Version returns the version number for this ROOT class
func (cls *class) Version() int {
	return cls.version
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (cls *class) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &class{
			members: make([]Member, 0),
		}
		return reflect.ValueOf(o)
	}
	Factory.db["TClass"] = f
	Factory.db["*rootio.class"] = f
}

// check interfaces
var _ Object = (*class)(nil)
var _ Class = (*class)(nil)
var _ ROOTUnmarshaler = (*class)(nil)

// EOF
