/* Look for file IOCs on the local system
 * jvehent - ulfr - 2013
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


type IOC struct {
	Raw, Path, Check, Value string
}


// core structure used to check each individual file and store results
type FileCheck struct {
	checks map[string]map[string]map[string]bool
}


/* ParseIOC takes a single IOC in the format <path>:<check>=<value>
   eg. /usr/bin/vim:md5=8680f252cabb7f4752f8927ce0c6f9bd
   return: IOC structure
*/
func ParseIOC(raw_ioc string) (ioc IOC) {
	ioc.Raw = raw_ioc
	tmp := strings.Split(raw_ioc, ":")
	ioc.Path = tmp[0]
	tmp = strings.Split(tmp[1], "=")
	ioc.Check, ioc.Value = tmp[0], tmp[1]
	return
}

/* GetFilesFromPath Walks through the path recursively and test all entries,
   when a file is found, it is added to the map with the associated IOC
   return: nil on success, error on failure
*/
func GetFilesFromPath(ioc IOC, files map[string]FileCheck) error {
	err := filepath.Walk(ioc.Path,
			     func(file string, f os.FileInfo, err error) error {
				if err != nil { panic(err) }
				if !f.IsDir() {
					//files[file] = ioc.Check: {
					//	ioc.Value: {
					//		ioc.Raw: false
					//	}
					//}
				}
				return nil
			     })
	if err != nil { panic(err) }
	return nil
}

/* GetFileMD5 opens a file and read it by blocks of 64 bytes (512 bits)
   update a md5 sum with each block
   this method plays nice with large files
   return: hex encoded MD5 hash
*/
func GetFileMD5(fp string) (hash string){
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
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func main() {
	iocs := make(map[string]IOC)
	files := make(map[string]FileCheck)
	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		raw_ioc := flag.Arg(i)
		iocs[raw_ioc] = ParseIOC(raw_ioc)
		GetFilesFromPath(iocs[raw_ioc], files)
	}
	for ioc, I := range iocs {
		fmt.Println(ioc, I.Path, I.Check, I.Value)
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
