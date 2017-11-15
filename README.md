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
go run ./main.go
```

# Expected result
```sh
ctx will work
found value: Go
key not found: color
found value: Go2
```

# Reference
1. https://github.com/googlearchive/ioweb2016/blob/master/backend/cache.go
1. https://golang.org/src/context/context.go#L398
1. https://godoc.org/golang.org/x/sync/syncmap
1. https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
