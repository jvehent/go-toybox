/* Look for file IOCs on the local system
   jvehent - ulfr - 2013
*/
package main
import(
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/* Representation of a File IOC.
	* Raw is the raw IOC string received from the program arguments
	* Path is the file system path to inspect
	* Check is the type of check, such as md5, sha1, contains, named, ...
	* Value is the value of the check, such as a md5 hash
	* Result is a boolean set to True when the IOC is detected
*/
type FileIOC struct {
	Raw, Path, Check, Value string
	Result bool
}


/* FileCheck is a structure used to perform checks against a file.
	* IOCs is an array of FileIOC
	* Checks is an array of Check used to speed up the lookups
*/
type FileCheck struct {
	IOCs []FileIOC
}


/* ParseIOC parses an IOC from the command line into a FileIOC struct
   parameters:
	* raw_ioc is a string that contains the IOC from the command line in
	the format <path>:<check>=<value>
	eg. /usr/bin/vim:md5=8680f252cabb7f4752f8927ce0c6f9bd
   return:
	* a FileIOC structure
*/
func ParseIOC(raw_ioc string) (ioc FileIOC) {
	ioc.Raw = raw_ioc
	tmp := strings.Split(raw_ioc, ":")
	ioc.Path = tmp[0]
	tmp = strings.Split(tmp[1], "=")
	ioc.Check, ioc.Value = tmp[0], tmp[1]
	return
}


/* BuildIOCChecklist takes an FileIOC structure, and walks through the path
   to list all the files that need to be inspected. When a file is found, it
   store it into the checklist map, with the associated FileIOC.
   parameters:
	* ioc is a FileIOC that contains a path to walk through
	* checklist is a checklist map to populate
   returns:
	* nil on success, error otherwise
*/
func BuildIOCChecklist(ioc FileIOC, checklist map[string]FileCheck) error {
	err := filepath.Walk(ioc.Path,
			     func(file string, f os.FileInfo, err error) error {
				if err != nil { panic(err) }
				if !f.IsDir() {
					var tmp = checklist[file]
					tmp.IOCs = append(tmp.IOCs, ioc)
					checklist[file] = tmp
				}
				return nil
			     })
	if err != nil { panic(err) }
	return nil
}

/* GetFileMD5 calculates the MD5 hash of a file.
   It opens a file, reads it by blocks of 64 bytes (512 bits), and updates a
   md5 sum with each block. This method plays nice with big files
   parameters:
	* fp is a string that contains the path of a file
   return:
	* hexhash, the hex encoded MD5 hash of the file found at fp
*/
func GetFileMD5(fp string) (hexhash string){
	h := md5.New()
	fd, err := os.Open(fp)
	if err != nil { panic(err)}
	defer func() {
		if err := fd.Close(); err != nil {
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
	hexhash = fmt.Sprintf("%x", h.Sum(nil))
	return
}


func main() {
	iocs := make(map[string]FileIOC)
	checklist := make(map[string]FileCheck)
	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		fmt.Println("Parsing IOC from command line", flag.Arg(i))
		raw_ioc := flag.Arg(i)
		iocs[raw_ioc] = ParseIOC(raw_ioc)
		if !BuildIOCChecklist(iocs[raw_ioc], checklist) {
			panic("Failed to build the checklist for ioc",raw_ioc)
		}
	}
	for f := range checklist{
		fmt.Println(f, checklist[f])
	}
	for _, I := range iocs {
		switch I.Check {
		case "contains":
		case "named":
		case "md5":
		case "sha1":
		case "sha256":
		case "sha512":
		case "sha3":
		default:
			err := "Unknown check type in " + I.Check
			panic(err)
		}
	}
	file := "/etc/passwd"
	hash := GetFileMD5(file)
	fmt.Println(hash, file)

	fd, err := os.Open(file)
	if err != nil { panic(err)}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(fd)
	re := regexp.MustCompile("ulfr")
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
