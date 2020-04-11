# si
Single Instance: keep just one instance of your program.

## Description
si is a Go(Lang) package that helps you to find and communicate with other instances of your program.

## Example
```go
package main

import (
	"fmt"
	"github.com/Miguel-Dorta/si"
	"os"
)

func main() {
	// Example 1: Register your instance and exit if there's another one running
	if err := si.Register("my awesome program"); err != nil {
		if err == si.ErrOtherInstanceRunning {
			fmt.Println("There's another instance running, bye!")
			os.Exit(0)
		}
		fmt.Println("Error registering instance: " + err.Error())
		os.Exit(1)
	}

	// Example 2: Find another instance and kill it
	p, err := si.Find("my awesome program")
	if err != nil {
		fmt.Println("Error finding another instance: " + err.Error())
		os.Exit(1)
	}
	if p == nil {
		fmt.Println("No other instance was found, bye!")
		os.Exit(0)
	}
	p.process.Kill()
}
```

### Important
If you want to communicate with your program using si.Process.StdXPipe(), this should have a pipe asigned to its stdX, else it will have /dev/null as default. This makes impossible the communication between processes. See [si/test/test.go/createInstance2()](https://github.com/Miguel-Dorta/si/blob/1689b65dcb56c7849fb21539e922639b6a76ded9/test/test.go#L61) to see an example.

#### Meme
![meme of a mexican dog saying "si"](https://i.nth.sh/media/4mM4w8KV46/DZCVdlgff1.jpeg "si")
