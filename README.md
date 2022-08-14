# goresult

## Install
```shell
go get github.com/LimitR/goresult
```

# Using
```go
package main

import (
	"github.com/LimitR/goresult"
)


func main() {
	c := goresult.NewResult("value ok")
	c2 := getResultOk()
	c3 := getResultError()
	fmt.Println(c.Unwrap()) // value ok
	fmt.Println(c2.Unwrap()) // ok
	fmt.Println(c3.Unwrap()) // 2022/08/14 16:25:47 Not ok
}

func getResultOk() *goresult.Result[any] {
	res := goresult.NewResult("ok")
	return res
}

func getResultError() *goresult.Result[any] {
	res := goresult.NewResult("ok")
	res.Error("Not ok")
	return res
}
```
## Check Error
```go
c3 := getResultError()
	if !c3.IsOk() {
		c3 = c3.Some("default")
	}
fmt.Println(c3.Unwrap()) // default

// Or
c3 := getResultError()

fmt.Println(c3.UnwrapOrElse("default")) // default
```