package main

import (
	"bytes"
	"fmt"
	"mig/modules/file"
)

func main() {
	var r file.Runner
	fmt.Println(r.Run(bytes.NewBuffer([]byte(`{"class":"parameters","parameters":{"searches":{"s1":{"paths":["/etc/passwd"],"contents":["ulfr"]}}}}`))))
}
