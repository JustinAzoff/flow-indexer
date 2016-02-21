//This shows some size comparisions of different encoding/compression methods
package main

import (
	"bytes"
	"compress/zlib"
	"fmt"

	"github.com/golang/snappy"
	"github.com/willf/bitset"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func delta_encode(docs []uint) []uint {
	encoded := make([]uint, len(docs))
	var last uint
	for i, val := range docs {
		encoded[i] = val - last
		last = val
	}
	return encoded
}

func makeBitset(docs []uint) []byte {
	bs := bitset.New(8)
	for _, id := range docs {
		bs.Set(id)
	}
	buffer := bytes.NewBuffer(make([]byte, 0, bs.BinaryStorageSize()))
	_, err := bs.WriteTo(buffer)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func makeMsgPack(docs []uint) []byte {
	b, err := msgpack.Marshal(delta_encode(docs))
	if err != nil {
		panic(err)
	}
	return b
}

func compressSnappy(buf []byte) []byte {
	compressed := snappy.Encode(nil, buf)
	return compressed
}

func compressZlib(buf []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(buf)
	w.Close()
	return b.Bytes()
}

func intRange(max int) []uint {
	nums := make([]uint, max)
	for i := 0; i < max; i++ {
		nums[i] = uint(i)
	}
	return nums
}

var encodingTests = []struct {
	description string
	ids         []uint
}{
	{"just 0", []uint{0}},
	{"just 1", []uint{1}},
	{"just 13", []uint{13}},
	{"just 23", []uint{23}},
	{"just 600", []uint{600}},
	{"just 6000", []uint{6000}},
	{"1,4,8,12", []uint{1, 4, 8, 12}},
	{"0-23", intRange(24)},
	{"0-720", intRange(720)},
	{"0-8640", intRange(8640)},
	{"0,1024", []uint{0, 1024}},
	{"rand", []uint{176, 370, 1138, 1308, 2435, 2441, 2559, 2621, 2646, 2870}},
	{"fib", []uint{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610}},
	{"fibmore", []uint{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711, 28657, 46368, 75025}},
}

func main() {
	fmt.Printf("%-10s: %10s %10s %10s %10s %10s %10s %10s %10s\n",
		"test",
		"bitset", "bitset.s",
		"msgpack", "msgpack.s",
		"msgpackd", "msgpackd.s", "msgpackd.z",
		"msgpackd.s/bitset.s",
	)
	for _, t := range encodingTests {
		bs := makeBitset(t.ids)
		lenBitset := len(bs)
		lenBitsetSnappy := len(compressSnappy(bs))

		mp := makeMsgPack(t.ids)
		lenMsgpack := len(mp)
		lenMsgpacksnappy := len(compressSnappy(mp))

		deltas := delta_encode(t.ids)
		mpd := makeMsgPack(deltas)
		lenMsgpackDelta := len(mpd)
		lenMsgpackDeltasnappy := len(compressSnappy(mpd))
		lenMsgpackDeltazlib := len(compressZlib(mpd))

		fmt.Printf("%-10s: %10d %10d %10d %10d %10d %10d %10d %10d%%\n",
			t.description,
			lenBitset, lenBitsetSnappy,
			lenMsgpack, lenMsgpacksnappy,
			lenMsgpackDelta, lenMsgpackDeltasnappy, lenMsgpackDeltazlib,
			uint64(100*float64(lenMsgpackDeltasnappy)/float64(lenBitsetSnappy)),
		)
		//fmt.Printf("\n%#v\nmessage pack delta %#v\nbitset: %#v\n", t.ids, mpd, bs)
	}
}
