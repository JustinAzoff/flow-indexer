package store

import (
	"reflect"
	"testing"
)

func intRange(max int) DocumentList {
	nums := make(DocumentList, max)
	for i := 0; i < max; i++ {
		nums[i] = DocumentID(i)
	}
	return nums
}

func runCodecTest(t *testing.T, codecFactory func() Codec) {
	ids := DocumentList{1, 2, 3, 5, 8, 13, 21}

	c := codecFactory()
	for _, id := range ids {
		b, _ := c.Bytes()
		c.FromBytes(b)
		c.AddID(id)
	}
	resultDocs := c.Documents()
	if !reflect.DeepEqual(ids, resultDocs) {
		t.Errorf("codec(%s)=> %v, want %v", c, resultDocs, ids)
	}

}

func runCodecBench(b *testing.B, codecFactory func() Codec, max int) {
	ids := intRange(max)
	for n := 0; n < b.N; n++ {
		c := codecFactory()
		c.AddID(0)
		for _, id := range ids {
			b, _ := c.Bytes()
			c.FromBytes(b)
			c.AddID(id)
		}
		c.Bytes()
	}
}

func TestCodec(t *testing.T) {
	//runCodecTest(t, NewBitsetCodec)
	//runCodecTest(t, NewMsgpackCodec)
	runCodecTest(t, func() Codec { return NewBitsetCodec() })
	runCodecTest(t, func() Codec { return NewMsgpackCodec() })
	runCodecTest(t, func() Codec { return NewMsgpackDeltasCodec() })
}

func BenchmarkBitsetCodec24(b *testing.B) {
	runCodecBench(b, func() Codec { return NewBitsetCodec() }, 24)
}
func BenchmarkMsgpackCodec24(b *testing.B) {
	runCodecBench(b, func() Codec { return NewMsgpackCodec() }, 24)
}

func BenchmarkMsgpackDeltaCodec24(b *testing.B) {
	runCodecBench(b, func() Codec { return NewMsgpackDeltasCodec() }, 24)
}

func BenchmarkBitsetCodec720(b *testing.B) {
	runCodecBench(b, func() Codec { return NewBitsetCodec() }, 720)
}
func BenchmarkMsgpackCodec720(b *testing.B) {
	runCodecBench(b, func() Codec { return NewMsgpackCodec() }, 720)
}

func BenchmarkMsgpackDeltaCodec720(b *testing.B) {
	runCodecBench(b, func() Codec { return NewMsgpackDeltasCodec() }, 720)
}
