package main

import (
	"context"
	"fmt"
)

var (
	ctx0 context.Context
)

type favContextKey string

func main() {
	initCtx()
	smallApplicationLevelCache()
}
func lookupCtx(c context.Context, k favContextKey) {
	if v := c.Value(k); v != nil {
		fmt.Println("found value:", v)
		return
	}
	fmt.Println("key not found:", k)
}

func initCtx() {
	ctx0 = context.Background()
}

func smallApplicationLevelCache() {
	k := favContextKey("language")
	k2 := favContextKey("language2")
	ctx0 = context.WithValue(ctx0, k, "Go")
	ctx0 = context.WithValue(ctx0, k2, "Go2")
	fmt.Println("ctx will work")
	lookupCtx(ctx0, k)
	lookupCtx(ctx0, favContextKey("color"))
	lookupCtx(ctx0, k2)
}
