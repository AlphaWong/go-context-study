package main

import (
	"context"
	"sync"
	"testing"
)

var (
	m   sync.Map
	ctx context.Context = context.Background()
)

func setSyncMap(sm *sync.Map, k string, v interface{}) {
	sm.Store(k, v)
}

func getSyncMap(sm *sync.Map, k string) interface{} {
	v, _ := sm.Load(k)
	return v
}

func setCtx(c context.Context, k string, v interface{}) {
	ctx = context.WithValue(c, k, v)
}

func getCtx(c context.Context, k string) interface{} {
	return c.Value(k)
}

var j int = 0

func BenchmarkTestingCtxSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		setCtx(ctx, string(j), j)
		j++
	}
}

func BenchmarkTestingCtxGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getCtx(ctx, string(j))
	}
}

var i int = 0

func BenchmarkTestingSyncMapSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		setSyncMap(&m, string(i), i)
		i++
	}
}

func BenchmarkTestingSyncMapGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getSyncMap(&m, string(i))
	}
}
