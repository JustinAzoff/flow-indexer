package store

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *DocumentID) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var tmp uint
		tmp, err = dc.ReadUint()
		(*z) = DocumentID(tmp)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DocumentID) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint(uint(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DocumentID) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint(o, uint(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DocumentID) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var tmp uint
		tmp, bts, err = msgp.ReadUintBytes(bts)
		(*z) = DocumentID(tmp)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

func (z DocumentID) Msgsize() (s int) {
	s = msgp.UintSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DocumentList) DecodeMsg(dc *msgp.Reader) (err error) {
	var xsz uint32
	xsz, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(xsz) {
		(*z) = (*z)[:xsz]
	} else {
		(*z) = make(DocumentList, xsz)
	}
	for bzg := range *z {
		{
			var tmp uint
			tmp, err = dc.ReadUint()
			(*z)[bzg] = DocumentID(tmp)
		}
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DocumentList) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for bai := range z {
		err = en.WriteUint(uint(z[bai]))
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DocumentList) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for bai := range z {
		o = msgp.AppendUint(o, uint(z[bai]))
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DocumentList) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var xsz uint32
	xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(xsz) {
		(*z) = (*z)[:xsz]
	} else {
		(*z) = make(DocumentList, xsz)
	}
	for cmr := range *z {
		{
			var tmp uint
			tmp, bts, err = msgp.ReadUintBytes(bts)
			(*z)[cmr] = DocumentID(tmp)
		}
		if err != nil {
			return
		}
	}
	o = bts
	return
}

func (z DocumentList) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (len(z) * (msgp.UintSize))
	return
}
