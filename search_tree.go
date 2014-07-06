package main

import (
	"fmt"
	"strings"
)

type STree struct {
	And, Or    []*STree
	Key, Value string
	Result     bool
}

var i int

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
		`name=bob OR name=alice OR heartbeattime>last{10m} AND (name=eve AND version=123) OR (name=zob AND heartbeattime>last{5m}) AND actionid=456`,
	}
	for _, search := range searches {
		search = strings.ToUpper(search)
		fmt.Printf("Search string: %s\n", search)
		root := new(STree)
		root = buildSearchTree(search)
		fmt.Println(walkTree(root, 0, ""))
		fmt.Printf("\n-------------------\n")
	}
}

func walkTree(t *STree, lvl int, cnj string) string {
	prefix := ""
	for i := 1; i <= lvl; i++ {
		prefix += "\t"
	}
	res := ""
	lvl++
	for _, leaf := range t.And {
		res += walkTree(leaf, lvl, " and")
	}
	for _, leaf := range t.Or {
		res += walkTree(leaf, lvl, " or")
	}
	if len(t.And) == 0 && len(t.Or) == 0 {
		res = "\n" + prefix + t.Value + cnj
	}
	return res
}

func buildSearchTree(search string) *STree {
	root := new(STree)
	branch := root
	fmt.Printf("\nbuilding tree for string: %s\n", search)
	encl_lvl := 0
	// parse search string char by char from left to right
	for pos := 0; pos < len(search); pos++ {
		char := string(search[pos])
		// skip substrings with parenthesis
		if char == "(" {
			encl_lvl++
			continue
		}
		if encl_lvl > 0 {
			if char == ")" {
				encl_lvl--
			}
			continue
		}
		// look for a space followed by a conjunction
		if char == " " {
			// make sure we have enough characters
			// to not get past the end of the string
			if pos+5 >= len(search) {
				break
			}
			if search[pos:pos+5] == " AND " {
				// process left part of the string
				str := getSubstring(search[0:pos])
				subbranch := buildSearchTree(str)
				branch.And = append(branch.And, subbranch)
				// process right part of the string
				str = getSubstring(search[pos+5 : len(search)])
				subbranch = buildSearchTree(str)
				branch.And = append(branch.And, subbranch)
				break
			} else if search[pos:pos+4] == " OR " {
				// process left part of the string
				str := getSubstring(search[0:pos])
				subbranch := buildSearchTree(str)
				branch.Or = append(branch.Or, subbranch)
				// process right part of the string
				str = getSubstring(search[pos+4 : len(search)])
				subbranch = buildSearchTree(str)
				branch.Or = append(branch.Or, subbranch)
				break
			}
		}
	}
	if branch.Value == "" {
		branch.Value = search[0:len(search)]
	}
	fmt.Printf("returning branch with value %s\n", branch.Value)
	return root
}

// remove leading and trailing spaces
// remove parenthesis if only one level of enclosure
// example: '   (something=something AND a=b)   '
// becomes: 'something=something AND a=b'
func getSubstring(str string) string {
	fmt.Println("getsubstring from", str)
	strlen := len(str)
	start := 0
	end := strlen
	openencl := 0
	closeencl := 0
	for pos := 0; pos <= strlen-1; pos++ {
		switch string(str[pos]) {
		case " ":
			start = pos + 1
		case "(":
			openencl++
			pos = strlen
		default:
			pos = strlen
		}
	}
	for pos := strlen - 1; pos >= 0; pos-- {
		switch string(str[pos]) {
		case " ":
			end = pos
		case ")":
			closeencl++
			pos = 0
		default:
			pos = 0
		}
	}
	if openencl == closeencl && openencl == 1 {
		start++
		end--
	}
	return str[start:end]
}
