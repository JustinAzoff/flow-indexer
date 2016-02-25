package store

import (
	"bytes"
	"io"

	"github.com/willf/bitset"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Codec interface {
	AddID(documentID int) error
	ReadFrom(io.Reader) error
	WriteTo(io.Writer) error
	FromBytes([]byte) error
	Bytes() ([]byte, error)
	Documents() []int
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

func (c *BitsetCodec) AddID(documentID int) error {
	c.bs.Set(uint(documentID))
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

func (c *BitsetCodec) Documents() []int {
	var docs []int
	for i, e := c.bs.NextSet(0); e; i, e = c.bs.NextSet(i + 1) {
		docs = append(docs, int(i))
	}
	return docs
}

type MsgpackCodec struct {
	buffer *[]byte
	ints   []int
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
	var ints []int
	err := msgpack.Unmarshal(b, &ints)
	c.ints = ints
	return err
}

func (c *MsgpackCodec) AddID(documentID int) error {
	c.ints = append(c.ints, documentID)
	return nil
}

func (c *MsgpackCodec) WriteTo(w io.Writer) error {
	return msgpack.NewEncoder(w).Encode(c.ints)
}
func (c *MsgpackCodec) Bytes() ([]byte, error) {
	b, err := msgpack.Marshal(c.ints)
	return b, err
}

func (c *MsgpackCodec) Documents() []int {
	return c.ints
}
