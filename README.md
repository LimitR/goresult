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
	fmt.Println(c3.Expect("Castom panic")) // panic: Castom panic

    file := goresult.CreateResultFrom(os.Open("/path/to/file.txt")).Unwrap()
	defer file.Close()
}

func getResultOk() *goresult.Result[string] {
	res := goresult.NewResult("ok")
	return res
}

func getResultError() *goresult.Result[string] {
	res := goresult.NewResult("ok")
	res.AddError("Not ok")
	return res
}
```
## Check Error
```go
c3 := getResultError()
	if !c3.IsOk() {
		c3 = goresult.NewResult("default")
	}
fmt.Println(c3.Unwrap()) // default

// Or

c3 := getResultError()

fmt.Println(c3.UnwrapOrElse("default")) // default
```

If the result is an error, the value will be deleted after processing

## Error Handling

```go
c3 := getResultError().UnwrapOrOn(func(res error) string {
	fmt.Println(res.Error()) // Not ok
	return "default"
})
fmt.Println(c3) // default

// Or

c3 := getResultError().UnwrapDelay(func(res string) {
	fmt.Println(res) // ok
	// ... Some code before the panic
	// ...
}) // panic: AAAA
```