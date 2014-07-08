// $ go run filepath_match.go
// File '/etc/cron.d' matches pattern '/etc/cron*'
// File '/etc/cron.d/0hourly' matches pattern '/etc/cron.*/*'
// File '/etc/cron.d/raid-check' matches pattern '/etc/cron.*/*'
// File '/etc/cron.d/unbound-anchor' matches pattern '/etc/cron.*/*'
// File '/etc/cron.daily' matches pattern '/etc/cron*'
// File '/etc/cron.daily/google-chrome' matches pattern '/etc/cron.*/*'
// File '/etc/cron.daily/google-talkplugin' matches pattern '/etc/cron.*/*'
// File '/etc/cron.daily/logrotate' matches pattern '/etc/cron.*/*'
// File '/etc/cron.daily/man-db.cron' matches pattern '/etc/cron.*/*'
// File '/etc/cron.daily/mlocate' matches pattern '/etc/cron.*/*'
// File '/etc/cron.deny' matches pattern '/etc/cron*'
// File '/etc/cron.hourly' matches pattern '/etc/cron*'
// File '/etc/cron.hourly/0anacron' matches pattern '/etc/cron.*/*'
// File '/etc/cron.hourly/mcelog.cron' matches pattern '/etc/cron.*/*'
// File '/etc/cron.monthly' matches pattern '/etc/cron*'
// File '/etc/cron.weekly' matches pattern '/etc/cron*'
// File '/etc/crontab' matches pattern '/etc/cron*'
// File '/etc/passwd' matches pattern '/etc/pass*'
// File '/etc/passwd-' matches pattern '/etc/pass*'
// File '/etc/passwdqc.conf' matches pattern '/etc/pass*'

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var patterns = [...]string{
		`/var/log/boot*`,
		`/var/cache/*/saved*`,
		`/etc/cron.*/*`,
		`/usr/*/vim`,
		`/sbin/*tables`,
	}

	for _, pattern := range patterns {
		// find the root path before the first pattern character.
		// seppos records the position of the latest path separator
		// before the first pattern.
		seppos := 1
		for cursor := 0; cursor < len(pattern); cursor++ {
			char := pattern[cursor]
			switch char {
			case '*', '?', '[', '{':
				// found pattern character. but ignore it if preceded by backslash
				if cursor > 0 {
					if pattern[cursor-1] == '\\' {
						break
					}
				}
				// exit the loop
				goto walk
			case os.PathSeparator:
				if cursor > 0 {
					seppos = cursor
				}
			}
		}
	walk:
		root := pattern[0 : seppos+1]
		fmt.Printf("Starting directory walk at '%s'\n", root)
		err := walkDir(root, pattern)
		if err != nil {
			panic(err)
		}
	}
}

func walkDir(root, pattern string) (err error) {
	if !matchSubPattern(root, pattern) {
		return nil
	}
	dir, err := os.Open(root)
	dirContent, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}
	// loop over the content of the directory
	for _, DirEntry := range dirContent {
		// if not a file
		if !DirEntry.Mode().IsRegular() {
			// if is a dir, recursively go down the path
			if DirEntry.IsDir() {
				path := root
				// append path separator if missing
				if path[len(path)-1] != os.PathSeparator {
					path += string(os.PathSeparator)
				}
				path += DirEntry.Name()
				err = walkDir(path, pattern)
				if err != nil {
					panic(err)
				}
			}
			// ignore non file
			continue
		}
		filename := root + string(os.PathSeparator) + DirEntry.Name()
		match, err := filepath.Match(pattern, filename)
		if err != nil {
			panic(err)
		}
		if match {
			fmt.Printf("File '%s' matches pattern '%s'\n", filename, pattern)
		}

	}
	dir.Close()
	return
}

// matchSubPattern is an optimization that matches a pattern with the current
// depth of a directory.
// To make filtering more efficient, split the pattern at the PathSeparator
// level of the current path. If the current levels don't match, there's no
// need to continue further down this path
func matchSubPattern(root, pattern string) bool {
	rootdepth := 0
	for pos := 0; pos < len(root); pos++ {
		if root[pos] == os.PathSeparator {
			rootdepth++
		}
	}
	subpattern := pattern
	patterndepth := 0
	for pos := 0; pos < len(pattern); pos++ {
		if pattern[pos] == os.PathSeparator {
			patterndepth++
		}
		if patterndepth == rootdepth {
			// pattern reaches the same depth as root, we have two choices:
			// 1. pattern has a match in the current depth, in which case we
			//    use pattern as it is
			// 2. pattern has a match in a subdirectory, so we create a subpattern
			//    that only matches the current depth
			subpattern = pattern[0 : pos+1]
			if -1 != strings.Index(pattern[pos+1:len(pattern)-1], string(os.PathSeparator)) {
				subpattern += "*"
			} else {
				subpattern = pattern
			}
			break
		}
	}
	match, _ := filepath.Match(subpattern, root)
	if !match {
		return false
	}
	return true
}
