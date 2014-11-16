Golang bindings for libvarnam
================

```golang
package main

import "fmt"
import "github.com/varnamproject/libvarnam-golang"

func main() {
	varnam, err := libvarnam.Init("hi")
	if err != nil {
		fmt.Errorf("Failed to initialize varnam. %s", err.Error())
	}

	words, err := varnam.Transliterate("bharat")
	if err != nil {
		fmt.Errorf("Failed to transliterate. %s", err.Error())
	}

	for _, word := range words {
		fmt.Printf("%s\n", word)
	}
}
```



