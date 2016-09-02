package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/willf/bitset"

	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/tecbot/gorocksdb"
)

type RocksDBStore struct {
	filename     string
	db           *gorocksdb.DB
	codecFactory func() Codec
	cfdocs       *gorocksdb.ColumnFamilyHandle
	cfips        *gorocksdb.ColumnFamilyHandle
}

var ro = gorocksdb.NewDefaultReadOptions()
var wo = gorocksdb.NewDefaultWriteOptions()

func NewRocksDBStore(filename string) (IpStore, error) {
	//TODO: steal options from ledisdb
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockSize(65536)
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.SetWriteBufferSize(134217728)
	db, cfh, err := gorocksdb.OpenDbColumnFamilies(opts, filename, []string{"default", "ips"}, []*gorocksdb.Options{opts, opts})
	if err != nil {
		return nil, errors.Wrap(err, "NewRocksDBStore failed")
	}
	newStore := &RocksDBStore{
		db:           db,
		filename:     filename,
		cfdocs:       cfh[0],
		cfips:        cfh[1],
		codecFactory: func() Codec { return NewBitsetCodec() },
	}
	return newStore, nil
}

func (rs *RocksDBStore) Close() error {
	rs.cfdocs.Destroy()
	rs.cfips.Destroy()
	rs.db.Close()
	return nil
}

func (rs *RocksDBStore) Compact() error {
	rs.db.CompactRange(gorocksdb.Range{nil, nil})
	return nil
}

func (rs *RocksDBStore) Filename() string {
	return rs.filename
}

func (rs *RocksDBStore) HasDocument(filename string) (bool, error) {
	key := []byte(filename)
	val, err := rs.db.GetCF(ro, rs.cfdocs, key)
	if err != nil {
		return false, errors.Wrap(err, "HasDocument")
	}
	defer val.Free()
	return val.Size() != 0, nil
}

func (rs *RocksDBStore) AddDocument(filename string, ips ipset.Set) error {
	exists, err := rs.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	nextID, err := rs.nextDocID()
	if err != nil {
		return err
	}

	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()
	rs.setDocId(batch, filename, nextID)

	for _, ip := range ips.SortedStrings() {
		//fmt.Printf("Add %#v to document\n", ip)
		err = rs.addIP(batch, nextID, ip)
		if err != nil {
			return err
		}
	}
	err = rs.db.Write(wo, batch)
	return err
}

func (rs *RocksDBStore) ListDocuments() error {
	nextID, err := rs.nextDocID()
	for i := uint64(0); i < nextID; i += 1 {
		name, err := rs.DocumentIDToName(i)
		if err != nil {
			break
		}
		fmt.Printf("Document %d is %#v\n", i, name)
	}
	return err
}

func (rs *RocksDBStore) DocumentIDToName(id uint64) (string, error) {
	idBytes := make([]byte, 10)
	binary.PutUvarint(idBytes[:], id)
	v, err := rs.db.GetCF(ro, rs.cfdocs, idBytes)
	defer v.Free()
	return string(v.Data()), err
}

func (rs *RocksDBStore) ExpandCIDR(ip string) ([]net.IP, error) {
	var ips []net.IP

	start, end, err := ipset.CIDRToByteStrings(ip)
	endBytes := []byte(end)
	if err != nil {
		return ips, err
	}

	it := rs.db.NewIteratorCF(ro, rs.cfips)
	defer it.Close()
	it.Seek([]byte(start))
	for it = it; it.Valid(); it.Next() {
		key := it.Key()
		if key.Size() != len(start) {
			//Ensure the matched keys are in the right ip family
			key.Free()
			continue
		}
		if bytes.Compare(key.Data(), endBytes) > 0 {
			key.Free()
			break
		}
		keycopy := make([]byte, key.Size())
		copy(keycopy, key.Data())
		key.Free()
		ip := net.IP(keycopy)
		ips = append(ips, ip)
	}
	err = it.Err()
	return ips, err
}

func (rs *RocksDBStore) QueryString(ip string) ([]string, error) {
	if strings.Contains(ip, "/") {
		return rs.QueryStringCidr(ip)
	}
	return rs.QueryStringIP(ip)
}

func (rs *RocksDBStore) QueryStringCidr(ip string) ([]string, error) {
	var start, end string
	start, end, err := ipset.CIDRToByteStrings(ip)
	endBytes := []byte(end)
	if err != nil {
		return nil, err
	}
	bs := bitset.New(8)

	codec := rs.codecFactory()

	it := rs.db.NewIteratorCF(ro, rs.cfips)
	defer it.Close()
	it.Seek([]byte(start))
	for it = it; it.Valid(); it.Next() {
		key := it.Key()
		if key.Size() != len(start) {
			//Ensure the matched keys are in the right ip family
			key.Free()
			continue
		}
		if bytes.Compare(key.Data(), endBytes) > 0 {
			key.Free()
			break
		}
		value := it.Value()
		codec.FromBytes(value.Data())
		value.Free()
		tmpBs := codec.ToBitset()
		bs = bs.Union(tmpBs)
	}
	err = it.Err()
	if err != nil {
		return nil, err
	}

	return rs.bitsetToDocs(bs)
}

func (rs *RocksDBStore) QueryStringIP(ip string) ([]string, error) {
	var docs []string
	key, err := ipset.IPStringToByteString(ip)
	if err != nil {
		return nil, err
	}
	v, err := rs.db.GetCF(ro, rs.cfips, []byte(key))
	defer v.Free()
	if v.Size() == 0 {
		return docs, nil
	}
	codec := rs.codecFactory()
	value := v.Data()
	codec.FromBytes(value)
	bs := codec.ToBitset()
	return rs.bitsetToDocs(bs)
}

func (rs *RocksDBStore) bitsetToDocs(bs *bitset.BitSet) ([]string, error) {
	var err error
	var docs []string
	for i, e := bs.NextSet(0); e; i, e = bs.NextSet(i + 1) {
		name, err := rs.DocumentIDToName(uint64(i))
		if err != nil {
			break
		}
		docs = append(docs, name)
	}
	return docs, err
}

func (rs *RocksDBStore) nextDocID() (uint64, error) {
	v, err := rs.db.GetCF(ro, rs.cfdocs, []byte("max_id"))
	if err != nil {
		return 0, errors.Wrap(err, "nextDocID")
	}
	defer v.Free()
	if v.Size() == 0 {
		return 0, nil
	}
	maxID, read := binary.Uvarint(v.Data())
	if read <= 0 {
		return 0, fmt.Errorf("nextDocID: Error converting %#v to a uint64", v)
	}
	return maxID + 1, nil

}

func (rs *RocksDBStore) setDocId(batch *gorocksdb.WriteBatch, filename string, id uint64) {
	idBytes := make([]byte, 10)
	binary.PutUvarint(idBytes[:], id)
	fnBytes := []byte(filename)
	batch.PutCF(rs.cfdocs, fnBytes, idBytes)
	batch.PutCF(rs.cfdocs, idBytes, fnBytes)
	batch.PutCF(rs.cfdocs, []byte("max_id"), idBytes)
}

func (rs *RocksDBStore) addIP(batch *gorocksdb.WriteBatch, id uint64, k string) error {
	v, err := rs.db.GetCF(ro, rs.cfips, []byte(k))
	if err != nil {
		return errors.Wrap(err, "addIP")
	}
	defer v.Free()
	codec := rs.codecFactory()
	if v.Size() != 0 {
		codec.FromBytes(v.Data())
	}
	codec.AddID(DocumentID(id))

	bytes, err := codec.Bytes()
	if err != nil {
		return err
	}
	batch.PutCF(rs.cfips, []byte(k), bytes)
	return nil
}

func init() {
	storeFactories["rocksdb"] = NewRocksDBStore
}
