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
	res.AddError(errors.New("Not ok"))
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
	return "default" // new value c3
})
fmt.Println(c3) // default

// Or

c3 := getResultError().UnwrapDelay(func(res string) {
	fmt.Println(res) // ok - value Result
	// ... Some code before the panic
	// ...
}) // panic: AAAA

// Or

c3 := getResultError().UnwrapDelay(func(res string) {
	fmt.Println(res) // ok - value Result
	recover()
}) // not panic
// But c3 = Result[string]
fmt.Println(c3.Unwrap()) // ok
```
If you need to process more than one Result
```go
c := goresult.NewResult("ok")
c2 := goresult.NewResult("ok2")
c2.AddError(errors.New("Panic"))
c3 := goresult.NewResult("ok3")

ch := []Result[string]{*c, *c2, *c3}

fmt.Println(goresult.CheckAll(ch)) // [ok ok3]
```
## Trace error

```go
func main() {
	result := a()
	value := result.Unwrap()
	fmt.Println(value)
}

func a() *goresult.Result[string] {
	result := b()
	// Some code...
	result.AddTrace()
	return result
}

func b() *goresult.Result[string] {
	result := c()
	// Some code...
	if !result.IsOk() {
		result.AddError(errors.New("Error in 'b'"))
	}
	return result
}

func c() *goresult.Result[string] {
	result := goresult.CreateResultFrom(d())
	// Some code...
	return result
}

func d() (string, error) {
	return "", errors.New("Error in 'd'")
}
```
Out:
```bash
panic:
Trace: main:12 -> a:19 -> b:27 -> c:33
Message: b 'Error in "b"', c 'Error in "d"'
```
## Disable trace
```go
// Some function
// ...
func c() *goresult.Result[string] {
	v, e := d()
	result := goresult.CreateResultFrom(v, e, false)
	// Some code...
	return result
}

func d() (string, error) {
	return "", errors.New("Error in 'd'")
}
```
Out:
```bash
panic: Error in 'b'
```
Output last error

# Get trace
```go
func main() {
	result := a()
	value := result.GetErrorTrace() // []trace
	fmt.Println(value[0].Message)  // Error in 'd'
	fmt.Println(value[0].FileName) // /home/.../main.go
	fmt.Println(value[0].FnName)   // c
	fmt.Println(value[0].Line)     // 37
}
```
# Context
```go
func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Millisecond)
	defer cancel()
	// Create channel for struct
	ch := make(chan *goresult.Result[string])
	// Started function in gorutine
	goresult.CreateResultCallback(ctx, d, ch)
	go func() {
		// Wait result or context cancellation
		fmt.Println(<-ch) // &Result{} - blank struct result
		wg.Done()
	}()
	wg.Wait()
}

func d() (string, error) {
	time.Sleep(5 * time.Second)
	return "", errors.New("Error in 'd'")
}

// But
func d() (string, error) {
	// time.Sleep(5 * time.Second)
	return "", errors.New("Error in 'd'")
}
//...
	fmt.Println(<-ch) // &Result{} - correct struct result
//...
```