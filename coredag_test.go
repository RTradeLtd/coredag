package coredag

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	pb "github.com/RTradeLtd/coredag/pb"
	proto "github.com/golang/protobuf/proto"
	"github.com/multiformats/go-multihash"
)

var (
	hashFunc = "sha2-256"
	ienc     = "json" // encoding type
	format   = "cbor" // serialization format
)

func Test_CoreDAG_Bad_Ienc_And_Format(t *testing.T) {
	type Data struct {
		Foo string
		Bar string
	}
	d := Data{"hello", "world"}
	dBytes, err := json.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}

	hashFuncCode, ok := multihash.Names[strings.ToLower(hashFunc)]
	if !ok {
		t.Fatal("bad hash func")
	}
	hashLength, ok := multihash.DefaultLengths[hashFuncCode]
	if !ok {
		t.Fatal("bad hash func code length")
	}
	// test a bad ienc
	if _, err := ParseInputs("jojo", format, bytes.NewReader(dBytes), hashFuncCode, hashLength); err == nil {
		t.Fatal("error expected")
	}
	// test a bad format
	if _, err := ParseInputs(ienc, "badformat", bytes.NewReader(dBytes), hashFuncCode, hashLength); err == nil {
		t.Fatal("error expected")
	}
}

func Test_CoreDAG_CBOR(t *testing.T) {
	hashFuncCode, ok := multihash.Names[strings.ToLower(hashFunc)]
	if !ok {
		t.Fatal("bad hash func")
	}
	hashLength, ok := multihash.DefaultLengths[hashFuncCode]
	if !ok {
		t.Fatal("bad hash func code length")
	}
	type args struct {
		ienc   string
		format string
	}
	tests := []struct {
		name string
		args args
	}{
		{"CBOR-CBOR", args{"cbor", "cbor"}},
		{"CBOR-DAG_CBOR", args{"cbor", "dag-cbor"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseInputs(
				tt.args.ienc,
				tt.args.format,
				bytes.NewReader([]byte("0")),
				hashFuncCode,
				hashLength,
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(nodes) == 0 {
				t.Fatal("bad node count")
			}
		})
	}
}

func Test_CoreDAG_JSON(t *testing.T) {
	type Data struct {
		Foo string
		Bar string
	}
	d := Data{"hello", "world"}
	dBytes, err := json.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	hashFuncCode, ok := multihash.Names[strings.ToLower(hashFunc)]
	if !ok {
		t.Fatal("bad hash func")
	}
	hashLength, ok := multihash.DefaultLengths[hashFuncCode]
	if !ok {
		t.Fatal("bad hash func code length")
	}
	type args struct {
		ienc   string
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"JSON-CBOR", args{"json", "cbor"}, false},
		{"JSON-DAG_CBOR", args{"json", "dag-cbor"}, false},
		{"JSON-PROTOBUF", args{"json", "protobuf"}, false},
		{"JSON-DAG_PB", args{"json", "dag-pb"}, false},
		{"Fail-Bad-IENC", args{"blahblah", "cbor"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseInputs(
				tt.args.ienc,
				tt.args.format,
				bytes.NewReader(dBytes),
				hashFuncCode,
				hashLength,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseInputs() err = %v, wantErr %v", err, tt.wantErr)
			}
			if len(nodes) == 0 && !tt.wantErr {
				t.Fatal("bad node count")
			}
		})
	}
}

func Test_CoreDAG_Protobuf(t *testing.T) {
	pbObject := pb.Data{Data: "hello"}
	pbBytes, err := proto.Marshal(&pbObject)
	if err != nil {
		t.Fatal(err)
	}
	hashFuncCode, ok := multihash.Names[strings.ToLower(hashFunc)]
	if !ok {
		t.Fatal("bad hash func")
	}
	hashLength, ok := multihash.DefaultLengths[hashFuncCode]
	if !ok {
		t.Fatal("bad hash func code length")
	}
	type args struct {
		ienc   string
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"PROTOBUF-PROTOBUF", args{"protobuf", "protobuf"}, false},
		{"PROTOBUF-DAG_PB", args{"protobuf", "dag-pb"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseInputs(
				tt.args.ienc,
				tt.args.format,
				bytes.NewReader(pbBytes),
				hashFuncCode,
				hashLength,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseInputs() err = %v, wantErr %v", err, tt.wantErr)
			}
			if len(nodes) == 0 && !tt.wantErr {
				t.Fatal("bad node count")
			}
		})
	}
}

func Test_CoreDAG_Raw(t *testing.T) {
	hashFuncCode, ok := multihash.Names[strings.ToLower(hashFunc)]
	if !ok {
		t.Fatal("bad hash func")
	}
	hashLength, ok := multihash.DefaultLengths[hashFuncCode]
	if !ok {
		t.Fatal("bad hash func code length")
	}
	type args struct {
		ienc   string
		format string
	}
	tests := []struct {
		name string
		args args
	}{
		{"RAW-CBOR", args{"raw", "cbor"}},
		{"RAW-DAG_CBOR", args{"raw", "dag-cbor"}},
		{"RAW-RAW", args{"raw", "raw"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseInputs(
				tt.args.ienc,
				tt.args.format,
				bytes.NewReader([]byte("0")),
				hashFuncCode,
				hashLength,
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(nodes) == 0 {
				t.Fatal("bad node count")
			}
		})
	}
	pbObject := pb.Data{Data: "hello"}
	pbBytes, err := proto.Marshal(&pbObject)
	if err != nil {
		t.Fatal(err)
	}
	tests = []struct {
		name string
		args args
	}{
		{"RAW-PROTOBUF", args{"raw", "protobuf"}},
		{"RAW-DAG_PB", args{"raw", "dag-pb"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ParseInputs(
				tt.args.ienc,
				tt.args.format,
				bytes.NewReader(pbBytes),
				hashFuncCode,
				hashLength,
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(nodes) == 0 {
				t.Fatal("bad node count")
			}
		})
	}
}
