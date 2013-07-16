/* Get all the MD5 from a dir
   jvehent - ulfr - 2013
*/
package main

import(
	"bufio"
	"code.google.com/p/go.crypto/sha3"
	"flag"
	"fmt"
	"io"
	"os"
)

var DEBUG bool = false
var VERBOSE bool = true

func GetFileSHA3(fp string) (hexhash string){
	if DEBUG {
		fmt.Printf("GetFileMD5: computing hash for '%s'\n", fp)
	}
	h := sha3.NewKeccak512()
	fd, err := os.Open(fp)
	if err != nil {
		fmt.Printf("GetFileMD5: can't get MD5 for %s: %s", fp, err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(fd)
	buf := make([]byte, 4096)
	for {
		block, err := reader.Read(buf)
		if err != nil && err != io.EOF { panic(err) }
		if block == 0 { break }
		h.Write(buf[:block])
	}
	hexhash = fmt.Sprintf("%x %s", h.Sum(nil), fp)
	return
}

func main() {
	if DEBUG { VERBOSE = true }

	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		if VERBOSE {
			fmt.Printf("using path'%s'\n", flag.Arg(i))
		}
		fmt.Println(GetFileSHA3(flag.Arg(i)))
	}
}
