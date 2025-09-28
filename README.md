# q.go

This is a Go function for print debugging, where all the printed statements go to a dedicated file.

The module exports a single function `Q()`, which you can pass any expression, and it will be logged to the file `/tmp/q.txt`.
This means you can log from any goroutine or process, and you can easily find any q-printed statements separate from the rest of your logging.

Here's a simple example:

```go
f, err := os.Stat("maybe_this_file_exists.txt")
q.Q(err)
```

If the first argument is a format string, then the string will be interpolated before it gets logged.
For example:

```go
name := "triangle"
sides := 3

q.Q("a %s has %d sides", name, sides)
// "a triangle has 3 sides"
```

As well as logging the value, `Q()` logs the name of the calling function, and the expression that you logged.
For example, `Q(2 + 2)` will be logged as `main: 2 + 2 = 4`.

## Example

Here's a longer example of a program that uses `Q()`:

```go
package main

import (
	"github.com/alexwlchan/q"
	"os"
)

func printShapeInfo(name string, sides int) {
	q.Q("a %s has %d sides", name, sides)
}

func main() {
	q.Q("hello world")

	q.Q(2 + 2)

	_, err := os.Stat("does_not_exist.txt")
	q.Q(err)

	printShapeInfo("triangle", 3)
}
```

And here's what gets written to `/tmp/q.txt`:

```console
$ cat /tmp/q.txt
main: "hello world"

main: 2 + 2 = 4

main: err = stat does_not_exist.txt: no such file or directory

printShapeInfo: a triangle has 3 sides
```
