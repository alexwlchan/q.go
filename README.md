# q.go

This is a Go function for print debugging, where all the printed statements go to a dedicated file.

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

```console
$ cat /tmp/q.txt
main: "hello world"

main: 2 + 2 = 4

main: err = stat does_not_exist.txt: no such file or directory

printShapeInfo: a triangle has 3 sides
```