/* Get all the MD5 from a dir
   jvehent - ulfr - 2013
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: " + os.Args[0] + " <file>")
	}
	fmt.Fprintf(os.Stderr, "computing sha3 hashes for '%s'\n", os.Args[1])
	h224 := sha3.New224()
	h256 := sha3.New256()
	h384 := sha3.New384()
	h512 := sha3.New512()
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(fd)
	buf := make([]byte, 4096)
	for {
		block, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if block == 0 {
			break
		}
		h224.Write(buf[:block])
		h256.Write(buf[:block])
		h384.Write(buf[:block])
		h512.Write(buf[:block])
	}
	fmt.Printf("sha224: %x %s\n", h224.Sum(nil), os.Args[1])
	fmt.Printf("sha256: %x %s\n", h256.Sum(nil), os.Args[1])
	fmt.Printf("sha384: %x %s\n", h384.Sum(nil), os.Args[1])
	fmt.Printf("sha512: %x %s\n", h512.Sum(nil), os.Args[1])
}
