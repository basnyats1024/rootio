package rootio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Decoder struct {
	buf *bytes.Buffer
	len int64
}

func NewDecoder(buf *bytes.Buffer) *Decoder {
	dec := &Decoder{
		buf: buf,
		len: int64(buf.Len()),
	}
	return dec
}

func NewDecoderFromBytes(data []byte) *Decoder {
	buf := make([]byte, len(data))
	copy(buf, data)
	return NewDecoder(bytes.NewBuffer(buf))
}

func NewDecoderFromReader(r io.Reader, size int) (*Decoder, error) {
	data := make([]byte, size)
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, fmt.Errorf("rootio.Decoder: read too few bytes [%v]. requested [%v]", n, size)
	}
	return NewDecoder(bytes.NewBuffer(data)), nil
}

func (dec *Decoder) Clone() *Decoder {
	o := NewDecoderFromBytes(dec.buf.Bytes())
	o.len = dec.len
	return o
}

func (dec *Decoder) Pos() int64 {
	return dec.len - int64(dec.buf.Len())
}

func (dec *Decoder) Len() int64 {
	return int64(dec.buf.Len())
}

func (dec *Decoder) readString(s *string) error {
	var err error
	var length byte
	var buf [256]byte

	err = dec.readBin(&length)
	if err != nil {
		return err
	}

	if length != 0 {
		err = dec.readBin(buf[:length])
		if err != nil {
			return err
		}
		*s = string(buf[:length])
	}
	return err

}

func (dec *Decoder) readCString(s *string, max int) error {
	var err error
	o := []byte{}
	n := 0
	var v byte
	for {
		err = dec.readBin(&v)
		if err != nil {
			return err
		}
		if v == 0 {
			break
		}
		n += 1
		if max > 0 && n >= max {
			break
		}
		o = append(o, v)
	}
	*s = string(o)
	return err
}

func (dec *Decoder) readBin(v interface{}) error {
	return binary.Read(dec.buf, E, v)
}

func (dec *Decoder) readInt16(v interface{}) error {
	var err error
	var d int16
	err = dec.readBin(&d)
	if err != nil {
		return err
	}

	switch uv := v.(type) {
	case *int16:
		*uv = int16(d)
	case *int32:
		*uv = int32(d)
	case *int64:
		*uv = int64(d)
	default:
		panic("Unknown type")
	}

	return err
}

func (dec *Decoder) readInt32(v interface{}) error {
	var err error
	switch uv := v.(type) {
	case *int32:
		err = dec.readBin(v)
	case *int64:
		var d int32
		err = dec.readBin(&d)
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
	return err
}

func (dec *Decoder) readInt64(v interface{}) error {
	var err error
	switch uv := v.(type) {
	case *int64:
		var d int64
		err = dec.readBin(&d)
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
	return err
}

func (dec *Decoder) readVersion() (version int16, position, bytecount int32, err error) {

	start := dec.Pos()

	tmp := dec.Clone()

	var bcnt uint32
	err = tmp.readBin(&bcnt)
	if err != nil {
		return
	}
	myprintf("readVersion - bytecount= %v\n", bcnt)
	if (int64(bcnt) & kByteCountMask) != 0 {
		bytecount = int32(int64(bcnt) & ^kByteCountMask)
		// as dec.buf was cloned, we need to read the 4-bytes holding the bytecount
		// we just read.
		dec.buf.Next(4)
	} else {
		// old version. no bytecount. next 4-bytes are holding the version.
		// as dec.buf has been left unchanged, we don't need to rewind to read them.
	}

	err = dec.readBin(&version)
	if err != nil {
		return
	}

	position = int32(start)
	myprintf("readVersion => [%v] [%v] [%v]\n", position, version, bytecount)
	return version, position, bytecount, err
}

func (dec *Decoder) readClass(name *string, count *int32, isref *bool) error {
	var err error
	*isref = false

	var tag uint32
	err = dec.readBin(&tag)
	if err != nil {
		return err
	}
	myprintf("::readClass. first int: [%v]\n", tag)
	switch {
	case tag == kNullTag:
		*isref = false
		*count = 0
		return err

	case (tag & kByteCountMask) != 0:
		// bufvers = 1
		classtag := ""
		err = dec.readClassTag(&classtag)
		if err != nil {
			return err
		}
		if classtag == "" {
			return fmt.Errorf("rootio.readClass: empty class tag")
		}
		*name = classtag
		*count = int32(int64(tag) & ^kByteCountMask)
		*isref = false
	default:
		*count = int32(tag)
		*isref = true
		*name = ""
	}
	return err
}

func (dec *Decoder) readClassTag(classtag *string) error {
	var err error
	var tag uint32
	err = dec.readBin(&tag)
	if err != nil {
		return err
	}

	tagNewClass := tag == kNewClassTag
	tagClassMask := (int64(tag) & (^int64(kClassMask))) != 0

	if tagNewClass {
		err = dec.readCString(classtag, 80)
		if err != nil {
			return err
		}
	} else if tagClassMask {
		panic("rootio.readClassTag: kClassMask not implemented")
	} else {
		panic(fmt.Errorf("rootio.readClassTag: unknown class-tag [%v]", tag))
	}

	return err
}

func (dec *Decoder) checkByteCount(pos, count int32, start int64, class string) error {
	if count == 0 {
		return nil
	}

	lenbuf := int64(pos) + int64(count) + 4
	diff := dec.Pos() - start
	if diff == lenbuf {
		return nil
	}
	err := fmt.Errorf(
		"**error** [%v] diff=%v len=%v (pos=%v, count=%v, start=%v)",
		class, diff, lenbuf, pos, count, start,
	)
	panic(err)
	return err
}

func (dec *Decoder) readObject(o *Object) error {
	//start := dec.Pos()
	//orig := dec.Clone()

	var class string
	var count int32
	var isref bool
	err := dec.readClass(&class, &count, &isref)
	if err != nil {
		return err
	}

	return err
}

// EOF
