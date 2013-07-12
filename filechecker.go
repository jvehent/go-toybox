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
		fmt.Printf("GetFileMD5: can't get MD5 for %s: %s", fp, err)
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
	* file is the absolute filename of the file to check
	* hash is the value of the hash being checked
	* check is the type of check
	* ActiveIOCIDs is a slice of int with IDs of active IOCs
	* IOCs is a map of IOC
   returns:
	* IsVerified: true if a match is found, false otherwise
*/
func VerifyHash(file string, hash string, check int, ActiveIOCIDs []int,
		IOCs map[int]FileIOC) (IsVerified bool) {
	IsVerified = false
	for _, id := range ActiveIOCIDs {
		if IOCs[id].Value == hash {
			IsVerified		= true
			tmpioc			:= IOCs[id]
			tmpioc.Result		= true
			tmpioc.ResultCount	+= 1
			tmpioc.Files		= append(tmpioc.Files, file)
			IOCs[id]		= tmpioc
		}
	}
	return
}


/* InspectFile is an orchestration function that runs the individual checks
   against a particular file. It uses the CheckBitMask to select which checks
   to run, and runs the checks in a smart way to minimize effort.
   parameters:
	* file is a string with the absolute path of the file that needs checking
	* ActiveIOCIDs is a slice of integer that contains the IDs of the IOCs
	that all files in that path and below must be checked against
	* CheckBitMask is a bitmask of the checks types currently active
	* IOCs is the global list of IOCs
   returns:
	* nil on success, error on failure
*/
func InspectFile(file string, ActiveIOCIDs []int, CheckBitMask int,
		 IOCs map[int]FileIOC) (error) {
	/* Iterate through the entire checklist, and process the checks of
	   each file
	*/
	if DEBUG {
		fmt.Printf("InspectFile: %s CheckMask %d\n", file, CheckBitMask)
	}
	if (CheckBitMask & CheckContains)	!= 0 {
		fmt.Println("Contains method not implemented")
	}
	if (CheckBitMask & CheckNamed)	!= 0 {
		fmt.Println("Contains method not implemented")
	}
	if (CheckBitMask & CheckMD5)		!= 0 {
		hash := GetHash(file, CheckMD5)
		if VerifyHash(file, hash, CheckMD5, ActiveIOCIDs, IOCs) {
			fmt.Printf("Positive result: %s\n", file)
		}
	}
	if (CheckBitMask & CheckSHA1)		!= 0 {
		fmt.Println("Contains method not implemented")
	}
	if (CheckBitMask & CheckSHA256)	!= 0 {
		fmt.Println("Contains method not implemented")
	}
	if (CheckBitMask & CheckSHA512)	!= 0 {
		fmt.Println("Contains method not implemented")
	}
	if (CheckBitMask & CheckSHA3)		!= 0 {
		fmt.Println("Contains method not implemented")
	}
	return nil
}

/* GetDownThatPath goes down a directory and build a list of Active IOCs that
   apply to the current path. For a given directory, it calls itself for all
   subdirectories fund, recursively walking down the pass. When it find a file,
   it calls the inspection function, and give it the list of IOCs to inspect
   the file with.
   parameters:
	* path is the file system path to inspect
	* ActiveIOCIDs is a slice of integer that contains the IDs of the IOCs
	that all files in that path and below must be checked against
	* CheckBitMask is a bitmask of the checks types currently active
	* IOCs is the global list of IOCs
	* ToDoIOCs is a map that contains the IOCs that are not yet active
   return:
	* nil on success, error on error
*/
func GetDownThatPath(path string, ActiveIOCIDs []int, CheckBitMask int,
		     IOCs map[int]FileIOC, ToDoIOCs map[int]FileIOC) (error) {
	for id, ioc := range ToDoIOCs {
		if ioc.Path == path {
			/* Found a new IOC to apply to the current path, add
			   it to the active list, and delete it from the todo
			*/
			ActiveIOCIDs = append(ActiveIOCIDs, id)
			CheckBitMask |= ioc.Check
			delete(ToDoIOCs, id)
		}
	}
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
			InspectFile(epath, ActiveIOCIDs, CheckBitMask, IOCs)
		}
	}
	for _, dir := range SubDirs {
		GetDownThatPath(dir, ActiveIOCIDs, CheckBitMask, IOCs, ToDoIOCs)
	}
	return nil
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

	// list of IOCs to process, remove from list when processed
	ToDoIOCs := make(map[int]FileIOC)

	flag.Parse()
	for i := 0; flag.Arg(i) != ""; i++ {
		if VERBOSE {
			fmt.Printf("Parsing IOC #%d '%s'\n", i, flag.Arg(i))
		}
		raw_ioc := flag.Arg(i)
		IOCs[i] = ParseIOC(raw_ioc, i)
		ToDoIOCs[i] = IOCs[i]
	}
	if VERBOSE {
		fmt.Println("Checklist built. Initiating inspection")
	}
	for id, ioc := range IOCs {
		// loop through the list of IOC, and only process the IOCs that
		// are still in the todo list
		if _, ok := ToDoIOCs[id]; !ok {
			continue
		}
		var EmptyActiveIOCs []int
		GetDownThatPath(ioc.Path, EmptyActiveIOCs, 0, IOCs, ToDoIOCs)
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
