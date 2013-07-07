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
	//"regexp"
	"strings"
)

var DEBUG bool = false
var VERBOSE bool = true

/* BitMask for the type of check to apply to a given file
   see documentation about iota for more info
*/
const(
	CheckContains	= 1 << iota
	CheckNamed
	CheckMD5
	CheckSHA1
	CheckSHA256
	CheckSHA512
	CheckSHA3
)


/* Representation of a File IOC.
	* Raw is the raw IOC string received from the program arguments
	* Path is the file system path to inspect
	* Value is the value of the check, such as a md5 hash
	* Check is the type of check in integer form
	* ResultCount is a counter of positive results for this IOC
	* Result is a boolean set to True when the IOC has matched once or more
	* Files is an slice of string that contains paths of matching files
*/
type FileIOC struct {
	Raw, Path, Value	string
	ID, Check, ResultCount	int
	Result			bool
	Files			[]string
}


/* FileCheck is a structure used to perform checks against a file
	* IOCs is a slice that contains the IDs of the IOCs to check.
	* Checks is a bitmask of said checks, for fast looking
*/
type FileCheck struct {
	IOCs		[]int
	CheckMask	int
}


/* ParseIOC parses an IOC from the command line into a FileIOC struct
   parameters:
	* raw_ioc is a string that contains the IOC from the command line in
	the format <path>:<check>=<value>
	eg. /usr/bin/vim:md5=8680f252cabb7f4752f8927ce0c6f9bd
	* id is an integer used as a ID reference
   return:
	* a FileIOC structure
*/
func ParseIOC(raw_ioc string, id int) (ioc FileIOC) {
	ioc.Raw		= raw_ioc
	ioc.ID		= id
	// split on the first ':' and use the left part as the Path
	tmp		:= strings.Split(raw_ioc, ":")
	ioc.Path	= tmp[0]
	// split the right part on '=', left is the check, right is the value
	tmp		= strings.Split(tmp[1], "=")
	ioc.Value	= tmp[1]
	// the check string is transformed into a bitmask and stored
	checkstring	:= tmp[0]
	switch checkstring {
	case "contains":
		ioc.Check = CheckContains
	case "named":
		ioc.Check = CheckNamed
	case "md5":
		ioc.Check = CheckMD5
	case "sha1":
		ioc.Check = CheckSHA1
	case "sha256":
		ioc.Check = CheckSHA256
	case "sha512":
		ioc.Check = CheckSHA512
	case "sha3":
		ioc.Check = CheckSHA3
	default:
		err := fmt.Sprintf("ParseIOC: Invalid check '%s'", checkstring)
		panic(err)
	}
	return
}


/* GetFilesFromPath walks through a path and lists all the files in contains
   parameters:
	* path, a string that contains the file system path to walk through
   returns:
	* Files, a string slice of files
*/
func GetFilesFromPath(path string) (Files []string) {
	err := filepath.Walk(path,
		func(file string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("GetFilesFromPath: error while",
					   "accessing %s: %s\n", file, err)
			}
			// Compare mode and perm bits to discard non-files
			fmode := f.Mode()
			if fmode.IsRegular() {
				Files = append(Files, file)
			}
			return nil
		})
	if err != nil { panic(err) }
	return Files
}


/* BuildIOCChecklist builds the IOC checklist for a given file
   parameters:
	* ioc is a FileIOC that contains a path to walk through
	* checklist is a checklist map to populate
	* file is the absolute path to a file
   returns:
	* nil on success, error otherwise
*/
func BuildIOCChecklist(	ioc FileIOC,
			Checklist map[string]FileCheck, file string) error {
	var chktmp = Checklist[file]
	chktmp.IOCs		= append(chktmp.IOCs, ioc.ID)
	chktmp.CheckMask	|= ioc.Check
	Checklist[file]		= chktmp
	return nil
}


/* GetHash is a wrapper above the hash functions
   parameters:
	* fp is a string that contains the path of a file
	* hash is an integer of the type of hash
   returns:
	* hexhash is a string that contains the hex encoded hash value
*/
func GetHash(fp string, hash int) (hexhash string) {
	switch hash{
	case CheckMD5:
		hexhash = GetFileMD5(fp)
	default:
		fmt.Printf("GetHash: '%s' is not a valid hash name\n", hash)
	}
	return
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
	if DEBUG {
		fmt.Printf("GetFileMD5: computing hash for '%s'\n", fp)
	}
	h := md5.New()
	fd, err := os.Open(fp)
	if err != nil {
		fmt.Printf("GetFileMD5: can get MD5 for %s: %s", fp, err)
	}
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


/* VerifyHash compares a file hash with the IOCs that apply to the file
   parameters:
	* Checklist is the global Checklist map. The map will be updated in
	  memory if a check is found.
	* IOCs is a map of IOC
	* Fp is the absolute filename of the file to check
	* HexHash is the value of the hash being checked
	* Check is the type of check
   returns:
	* IsVerified: true if a match is found, false otherwise
*/
func VerifyHash(Checklist map[string]FileCheck, IOCs map[int]FileIOC,
		Fp string, HexHash string, Check int) (IsVerified bool) {
	IsVerified = false
	for _, IOCID := range Checklist[Fp].IOCs {
		if IOCs[IOCID].Value == HexHash {
			IsVerified		= true
			tmpioc			:= IOCs[IOCID]
			tmpioc.Result		= true
			tmpioc.ResultCount	+= 1
			tmpioc.Files		= append(tmpioc.Files, Fp)
			IOCs[IOCID]		= tmpioc
		}
	}
	return
}


func main() {
	if DEBUG { VERBOSE = true }
	/* IOCs is a map of individual IOCs and associated results
		IOCs = {
			<id> = { <struct FileIOC> },
			<id> = { <struct FileIOC> },
			...
		}
	*/
	IOCs := make(map[int]FileIOC)

	/* Checklist contains a map of files and associated checks.
		Checklist = {
			'<file>' = {	checkmask: <bitmask>,
					IOCs: [<IocID>, ...]},
			'<file>' = {	<struct FileCheck> },
			'<file>' = {	<struct FileCheck> },
			...
		}
	*/
	Checklist := make(map[string]FileCheck)

	/* Files is a map that list all the files contained in a path. It is
	   used as a caching mechanism to reduce the number of times a path is
	   walked through when multiple IOCs reference the same path
	*/
	Files := make(map[string][]string)
	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		if VERBOSE {
			fmt.Printf("Parsing IOC #%d '%s'\n", i, flag.Arg(i))
		}
		raw_ioc := flag.Arg(i)
		IOCs[i] = ParseIOC(raw_ioc, i)
		if Files[IOCs[i].Path] == nil {
			Files[IOCs[i].Path] = GetFilesFromPath(IOCs[i].Path)
		}
		if VERBOSE {
			fmt.Printf("Loading %d files for IOC %s\n",
				   len(Files[IOCs[i].Path]), IOCs[i].Raw)
		}
		for _, file := range Files[IOCs[i].Path] {
			err := BuildIOCChecklist(IOCs[i], Checklist, file)
			if err != nil { panic(err) }
		}
	}
	if VERBOSE {
		fmt.Println("Checklist built. Initiating inspection")
	}
	if DEBUG {
		for file, check := range Checklist {
			fmt.Println(file, check)
		}
	}
	/* Iterate through the entire checklist, and process the checks of
	   each file
	*/
	for File, FileCheckList := range Checklist {
		if DEBUG {
			fmt.Printf("%s CheckMask %d\n", File,
				   FileCheckList.CheckMask)
		}
		if (FileCheckList.CheckMask & CheckContains)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		if (FileCheckList.CheckMask & CheckNamed)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		if (FileCheckList.CheckMask & CheckMD5)		!= 0 {
			FileHash := GetHash(File, CheckMD5)
			if VerifyHash(Checklist, IOCs, File, FileHash, CheckMD5) {
				fmt.Printf("Positive result: %s\n", File)
			}
		}
		if (FileCheckList.CheckMask & CheckSHA1)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		if (FileCheckList.CheckMask & CheckSHA256)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		if (FileCheckList.CheckMask & CheckSHA512)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		if (FileCheckList.CheckMask & CheckSHA3)	!= 0 {
			fmt.Println("Contains method not implemented")
		}
		// Done with this file, clean up
		delete(Checklist, File)
	}
	if VERBOSE {
		for _, ioc := range IOCs {
			fmt.Printf("IOC '%s' returned %d positive match\n",
				   ioc.Raw, ioc.ResultCount)
			if ioc.Result {
				for _, file := range ioc.Files {
					fmt.Printf("\t- %s\n", file)
				}
			}
		}
	}
	/*
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
	*/
}
