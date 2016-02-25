package store

import (
	"bytes"
	"io"

	"github.com/tinylib/msgp/msgp"
	"github.com/willf/bitset"
)

type Codec interface {
	AddID(DocumentID) error
	ReadFrom(io.Reader) error
	WriteTo(io.Writer) error
	FromBytes([]byte) error
	Bytes() ([]byte, error)
	Documents() DocumentList
	String() string
}

type BitsetCodec struct {
	buffer *[]byte
	bs     *bitset.BitSet
}

func NewBitsetCodec() *BitsetCodec {
	bs := bitset.New(8)
	return &BitsetCodec{bs: bs}
}

func (c *BitsetCodec) String() string {
	return "BitsetCodec"
}

func (c *BitsetCodec) ReadFrom(r io.Reader) error {
	c.bs.ReadFrom(r)
	return nil
}
func (c *BitsetCodec) FromBytes(b []byte) error {
	c.bs.ReadFrom(bytes.NewBuffer(b))
	return nil
}

func (c *BitsetCodec) AddID(id DocumentID) error {
	c.bs.Set(uint(id))
	return nil
}

func (c *BitsetCodec) WriteTo(w io.Writer) error {
	_, err := c.bs.WriteTo(w)
	return err
}
func (c *BitsetCodec) Bytes() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, c.bs.BinaryStorageSize()))
	err := c.WriteTo(buffer)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func (c *BitsetCodec) Documents() DocumentList {
	var docs DocumentList
	for i, e := c.bs.NextSet(0); e; i, e = c.bs.NextSet(i + 1) {
		docs = append(docs, DocumentID(i))
	}
	return docs
}

type MsgpackCodec struct {
	docs DocumentList
}

func NewMsgpackCodec() *MsgpackCodec {
	return &MsgpackCodec{}
}

func (c *MsgpackCodec) String() string {
	return "MsgpackCodec"
}

func (c *MsgpackCodec) ReadFrom(r io.Reader) error {
	return nil
}
func (c *MsgpackCodec) FromBytes(b []byte) error {
	_, err := c.docs.UnmarshalMsg(b)
	return err
}

func (c *MsgpackCodec) AddID(id DocumentID) error {
	c.docs = append(c.docs, id)
	return nil
}

func (c *MsgpackCodec) WriteTo(w io.Writer) error {
	mw := msgp.NewWriter(w)
	err := c.docs.EncodeMsg(mw)
	mw.Flush()
	return err
}
func (c *MsgpackCodec) Bytes() ([]byte, error) {
	b, err := c.docs.MarshalMsg([]byte{})
	return b, err
}

func (c *MsgpackCodec) Documents() DocumentList {
	return c.docs
}

type MsgpackDeltasCodec struct {
	docs        DocumentList
	encodedDocs DocumentList
}

func NewMsgpackDeltasCodec() *MsgpackDeltasCodec {
	return &MsgpackDeltasCodec{}
}

func (c *MsgpackDeltasCodec) String() string {
	return "MsgpackDeltasCodec"
}

func (c *MsgpackDeltasCodec) ReadFrom(r io.Reader) error {
	return nil
}
func (c *MsgpackDeltasCodec) FromBytes(b []byte) error {
	_, err := c.docs.UnmarshalMsg(b)
	if err != nil {
		return err
	}
	deltaDecode(c.docs)
	return nil
}

func (c *MsgpackDeltasCodec) AddID(id DocumentID) error {
	c.docs = append(c.docs, id)
	return nil
}

func (c *MsgpackDeltasCodec) WriteTo(w io.Writer) error {
	mw := msgp.NewWriter(w)
	c.encodedDocs = deltaEncode(c.docs)
	err := c.encodedDocs.EncodeMsg(mw)
	mw.Flush()
	return err
}
func (c *MsgpackDeltasCodec) Bytes() ([]byte, error) {
	c.encodedDocs = deltaEncode(c.docs)
	b, err := c.encodedDocs.MarshalMsg([]byte{})
	return b, err
}

func (c *MsgpackDeltasCodec) Documents() DocumentList {
	return c.docs
}

func deltaEncode(docs DocumentList) DocumentList {
	encoded := make(DocumentList, len(docs))
	var last DocumentID
	for i, val := range docs {
		encoded[i] = val - last
		last = val
	}
	return encoded
}
func deltaDecode(docs DocumentList) {
	var last DocumentID
	for i, val := range docs {
		docs[i] = val + last
		last = docs[i]
	}
}
