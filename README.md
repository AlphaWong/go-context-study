# Objective
I would like to test the application level cache.

# Issue background
I would like to store a TCP response in key value format. Those response not need to cleanup ( never expire until server restart ).  It needs to be thread safe.

# Original design
Create a public var as type sync.Map. It seems very dirty in the code level that I need to wrapper a get and set function and a not found function.

# After looking GoogleIO2016
A public var type ``Context`` can be considered as it is thread safe and It do return nil if the item is not found. 

# sync.Map
```go
# https://github.com/golang/sync/blob/master/syncmap/map.go#L98
func (m *Map) Load(key interface{}) (value interface{}, ok bool)
```

# Context
```go
# https://golang.org/src/context/context.go#L183
func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}
```

# Before
```go
// Load will return nil and false if key not found
func GetCacheAPIKeys(c string) (map[string]string, bool) {
  var emptyMap map[string]string
  key, ok := apiKeys.Load(c)
  if !ok{
    return emptyMap, ok
  }
  return key.(map[string]string), ok
}
```

# After
```go
// It will return nil if key not found
func GetCacheAPIKeys(c string) interface{} {
	return apiKeys.Value(c)
}
```

# Run
```sh
$ go test -bench=BenchmarkTesting -benchmem
goos: linux
goarch: amd64
pkg: github.com/go-context-study
BenchmarkTestingCtxSet-4       	 5000000	       222 ns/op	      79 B/op	       4 allocs/op
BenchmarkTestingCtxGet-4       	30000000	        82.7 ns/op	      20 B/op	       2 allocs/op
BenchmarkTestingSyncMapSet-4   	10000000	       143 ns/op	      48 B/op	       4 allocs/op
BenchmarkTestingSyncMapGet-4   	30000000	        39.3 ns/op	       0 B/op	       0 allocs/op
```

# Expected result
```sh
Ctx should always faster than sync Map
```

# Actual result
```
Sync map always faster than ctx
```

# Reason
```
sync map is hash based

ctx is a linked list
```

# ctx code
```go
// https://golang.org/src/context/context.go#L477<Paste>
// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.val)
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}
```

# Sum up
Both sync map and ctx is thread safe. However, ctx will always O(n) for get function.
Meanwhile sync map is hash based that get function will always O(1).

# Reference
1. https://github.com/googlearchive/ioweb2016/blob/master/backend/cache.go
1. https://golang.org/src/context/context.go#L398
1. https://godoc.org/golang.org/x/sync/syncmap
1. https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
