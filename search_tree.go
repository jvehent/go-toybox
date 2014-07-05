package main

import (
	"fmt"
	"strings"
)

type Stree struct {
	Left       *Stree
	Right      *Stree
	Key, Value string
}

func main() {
	var searches = [...]string{
		`name=bob`,
		`name=bob OR name=alice OR heartbeattime>last{5m}`,
		`name=bob AND version=123 AND heartbeattime>last{5m}`,
		`name=bob AND version=123 OR heartbeatttime>last{5m}`,
		`(name=bob AND version=123) OR (name=alice AND heartbeattime>last{5m}) AND actionid=456`,
		`(name=bob OR (name=alice AND version=234)) AND heartbeattime>last{5m}`,
		`version=123 AND (name=bob OR (name=alice AND version=234)) AND heartbeattime>last{5m} AND actionid=5`,
		`(name=alice OR (name=bob OR (name=eve AND version=1) AND version=2) AND version=3)`,
		`version=123 AND (name=alice OR (name=bob OR (name=eve AND version=1) AND version=2) AND version=3)`,
	}
	for _, search := range searches {
		search = strings.ToUpper(search)
		fmt.Printf("Parsing '%s'\n", search)
		//var tree *Stree
		// find non-enclosed AND conjunctions
		encl_lvl := 0
		encl_start_pos := 0
		var enclosures []string
		var andmarkers []int
		strlen := len(search)
		for pos := 0; pos < strlen; pos++ {
			char := string(search[pos])
			// look for parenthesis
			if char == "(" {
				if encl_lvl == 0 {
					encl_start_pos = pos
				}
				encl_lvl++
				continue
			}
			if encl_lvl > 0 {
				if char == ")" {
					encl_lvl--
					if encl_lvl == 0 {
						enclosures = append(enclosures, search[encl_start_pos+1:pos])
					}
				}
				continue
			}
			// look for a space followed by "AND"
			if char == " " {
				if len(search[pos:len(search)]) > 4 {
					if search[pos+1:pos+5] == "AND " {
						andmarkers = append(andmarkers, pos+1)
						// found, jump the next 3 char
						pos += 3
					}
				}
			}
		}
		fmt.Printf("Found AND markers at positions")
		for _, pos := range andmarkers {
			fmt.Printf(" %d;", pos)
		}
		fmt.Printf("\n")
		fmt.Printf("Enclosures: ")
		if len(enclosures) == 0 {
			fmt.Printf("None")
		} else {
			for i, encl := range enclosures {
				fmt.Printf("%d=(%s); ", i, encl)
			}
		}
		fmt.Printf("\n\n")
	}
}
