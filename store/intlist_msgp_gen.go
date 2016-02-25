package store

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Intlist) DecodeMsg(dc *msgp.Reader) (err error) {
	var xsz uint32
	xsz, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(xsz) {
		(*z) = (*z)[:xsz]
	} else {
		(*z) = make(Intlist, xsz)
	}
	for bzg := range *z {
		(*z)[bzg], err = dc.ReadInt()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Intlist) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for bai := range z {
		err = en.WriteInt(z[bai])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Intlist) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for bai := range z {
		o = msgp.AppendInt(o, z[bai])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Intlist) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var xsz uint32
	xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(xsz) {
		(*z) = (*z)[:xsz]
	} else {
		(*z) = make(Intlist, xsz)
	}
	for cmr := range *z {
		(*z)[cmr], bts, err = msgp.ReadIntBytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

func (z Intlist) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (len(z) * (msgp.IntSize))
	return
}
