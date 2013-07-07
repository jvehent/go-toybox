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
	* Check is the type of check, such as md5, sha1, contains, named, ...
	* Value is the value of the check, such as a md5 hash
	* Result is a boolean set to True when the IOC is detected
*/
type FileIOC struct {
	Raw, Path, Value	string
	ID, Check, ResultCount	int
	Result			bool
}


/* IOCCheck is the light version of an IOC struct used inside the checklist
	* ID maps to a FileIOC ID
	* Value maps to a FileIOC Value
	* Result is a boolean set to true when this particular check matches
*/
type IOCCheck struct {
	ID	int
	Value	string
	Result	bool
}


/* FileCheck is a structure used to perform checks against a file.
	* IOCs is an array of FileIOC
	* Checks is an array of unique check type used to speed up processing.
	If a given file needs to be checks for 3 md5 IOCs, the Checks array
	will contain ['md5'] once. The hash of the file will be computed once,
	and them compared to all 3 hashes.
*/
type FileCheck struct {
	IOCs		map[int][]IOCCheck
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
	ioc.Raw = raw_ioc
	ioc.ID = id
	// split on the first ':' and use the left part as the Path
	tmp := strings.Split(raw_ioc, ":")
	ioc.Path = tmp[0]
	// split the right part on '=', left is the check, right is the value
	tmp = strings.Split(tmp[1], "=")
	ioc.Value = tmp[1]
	// the check string is transformed into a bitmask and stored
	checkstring := tmp[0]
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
			if err != nil {
				fmt.Printf("BuildIOCChecklist: error while",
					"accessing %s: %s\n", file, err)
			}
			// Compare mode and perm bits to discard non-files
			fmode := f.Mode()
			if fmode.IsRegular() {
				/* We have a file. Add it to the checklist.
				   1. grab the pointer of the file entry into chktmp.
				   2. allocate the IOCs map if needed.
				   3. Update the CheckMask bitmask with current check
				   4. Allocate and populate the IOCCheck structure
				   5. Append the IOCCheck struct to the IOCs slice
				   6. Store the chktmp check into the checklist
				*/
				var chktmp = checklist[file]
				if chktmp.IOCs == nil {
					chktmp.IOCs = make(map[int][]IOCCheck)
				}
				chktmp.CheckMask |= ioc.Check
				var ioctmp IOCCheck
				ioctmp.ID = ioc.ID
				ioctmp.Value = ioc.Value
				ioctmp.Result = false
				chktmp.IOCs[ioc.Check] = append(chktmp.IOCs[ioc.Check], ioctmp)
				checklist[file] = chktmp
			}
			return nil
		})
	if err != nil { panic(err) }
	return nil
}


/* GetHash is a wrapper above the hash functions
   parameters:
	* fp is a string that contains the path of a file
	* hash is a string with the name of the hash to calculate
   returns:
	* hexhash is the hex encoded hash value
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
	for pos, ioc := range Checklist[Fp].IOCs[Check] {
		if ioc.Value == HexHash {
			IsVerified = true
			Checklist[Fp].IOCs[Check][pos].Result = true
			// store updated IOC results in IOCs list
			tmpioc := IOCs[ioc.ID]
			tmpioc.Result = true
			tmpioc.ResultCount++
			IOCs[ioc.ID] = tmpioc
		}
	}
	return
}


func main() {
	/* IOCs is a map of IOC received from the command line. each IOC receives
	   an integer ID used as a reference in the map.
	*/
	IOCs := make(map[int]FileIOC)

	/* Checklist is the core structure that represents IOC checks against files
	   Each individual file path has a FileCheck structure, that contains
	   one or more IOCCheck structure into a slices of checks.
	   The checks are grouped by type, so that all MD5 IOCChecks on a single
	   file will be listed inside 'Checklist[file].IOCs[CheckMD5]'
	   (CheckMD5, and all check types, are integer constants)
	Checklist: {
		"<path>(string)": {
			"CheckMask": (integer),					\
			"IOCs": {						|
				"<check>(integer)": [				|File
					{	'ID': (integer),	\ IOC	|
						'Value': (string),	| Check	|Check
						'Result': (bool)	/ struct|
					}					|struct
					{ ...}					|
				]						|
			}							|
		}								/
	}
	*/
	Checklist := make(map[string]FileCheck)
	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		fmt.Println("Parsing IOC from command line", flag.Arg(i))
		raw_ioc := flag.Arg(i)
		IOCs[i] = ParseIOC(raw_ioc, i)
		err := BuildIOCChecklist(IOCs[i], Checklist)
		if err != nil {
			panic(err)
		}
	}
	if DEBUG {
		for file, check := range Checklist {
			fmt.Println(file, check)
		}
		fmt.Printf("%d files to inspect\n", len(Checklist))
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
	}
	for _, ioc := range IOCs {
		fmt.Printf("IOC '%s' returned %d positive match\n",
			   ioc.Raw, ioc.ResultCount)
	}
	for fp, filechecks := range Checklist {
		for _, IOCslice := range filechecks.IOCs {
			for _, ioc := range IOCslice {
				if ioc.Result {
					fmt.Printf("'%s' positive to '%s'\n",
						   fp, IOCs[ioc.ID].Raw)
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
