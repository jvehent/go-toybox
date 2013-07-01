package main

import(
	"io"
	"os"
	"bufio"
	"fmt"
	"crypto/md5"
)

func get_md5_from_file(fp string) (string){
	/* read a file from its 'fp' path by block of 64 bytes (512 bits)
	*  to match MD5's internal block size, and play nice with big files
	*  returns: hex encoded MD5 hash
	*/
	h := md5.New()
	fd, err := os.Open(fp)
	if err != nil { panic(err)}
	defer func() {
		if err := fd.Close(); err !=nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(fd)
	buf := make([]byte, 64)
	for {
		block, err := reader.Read(buf)
		if err != nil && err != io.EOF { panic(err) }
		if block == 0 { break }
		h.Write(buf[:block])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
func main() {
	file := "/etc/passwd"
	fmt.Println("Calculating MD5 of", file)
	hash := get_md5_from_file(file)
	fmt.Println(hash, file)
}
