package main

import (
	"fmt"
	"log"

	"github.com/cloudfoundry-incubator/candiedyaml"
	//"gopkg.in/yaml.v2"
	"github.com/DirectXMan12/go-yaml"
)

var data = `
b:
  d: [3, 4]
  c: 2
# I wonder if this will stick
a: Easy!
`

func main() {
	fmt.Println("With GOYaml")
	var ms yaml.CommentedMapItem
	err := yaml.Unmarshal([]byte(data), &ms)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	data2, err := yaml.Marshal(ms)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s\n", data2)

	fmt.Println("With CandiedYaml")
	var s interface{}
	err = candiedyaml.Unmarshal([]byte(data), &s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	data2, err = candiedyaml.Marshal(s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s\n", data2)

}
