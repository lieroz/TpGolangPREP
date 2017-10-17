package main

import (
	"encoding/json"
	"testing"
	"github.com/json-iterator/go"
)

var (
	data = []byte(`{"RealName":"Vasily", "Login":"v.romanov", "Status":1, "Flags": 1}`)
	u    = User{}
	c    = Client{}
	c1   = Client{}
)

// go test -v -bench=. -benchmem json/*.go
// go test -v -bench=. json/*.go

func BenchmarkDecodeJsoniter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := jsoniter.ConfigFastest.BorrowIterator(data)
		iter.ReadVal(&c1)
		jsoniter.ConfigFastest.ReturnIterator(iter)
	}
}

func BenchmarkDecodeStandart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, &c)
	}
}

func BenchmarkDecodeEasyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = u.UnmarshalJSON(data)
	}
}

func BenchmarkEncodeJsoniter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stream := jsoniter.ConfigFastest.BorrowStream(nil)
		stream.WriteVal(c1)
		jsoniter.ConfigFastest.ReturnStream(stream)
	}
}

func BenchmarkEncodeStandart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&c)
	}
}

func BenchmarkEncodeEasyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalJSON()
	}
}
