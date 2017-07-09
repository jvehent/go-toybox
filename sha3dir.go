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
	h := sha3.NewKeccak224()
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


func GetDownThatPath(path string) (error) {
	var SubDirs []string
	/* Non-recursive directory walk-through. Read the content of dir stored
	   in 'path', put all sub-directories in the SubDirs slice, and call
	   the inspection function for all files
	*/
	cdir, err := os.Open(path)
	defer func() {
		if err := cdir.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil { panic(err) }
	dircontent, err := cdir.Readdir(-1)
	if err != nil { panic(err) }
	for _, entry := range dircontent {
		epath := path + "/" + entry.Name()
		if entry.IsDir() {
			SubDirs = append(SubDirs, epath)
		}
		if entry.Mode().IsRegular() {
			fmt.Println(GetFileSHA3(epath))
			//GetFileSHA3(epath)
		}
	}
	for _, dir := range SubDirs {
		GetDownThatPath(dir)
	}
	return nil
}


func main() {
	if DEBUG { VERBOSE = true }

	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		if VERBOSE {
			fmt.Printf("using path'%s'\n", flag.Arg(i))
		}
		GetDownThatPath(flag.Arg(i))
	}
}
