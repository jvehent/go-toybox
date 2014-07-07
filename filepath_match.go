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
)

func main() {
	var patterns []string
	patterns = append(patterns, `/etc/pass*`)
	patterns = append(patterns, `/etc/cron*`)
	patterns = append(patterns, `/etc/cron.*/*`)

	err := filepath.Walk("/etc", func(path string, _ os.FileInfo, _ error) error {
		err := evaluate(path, patterns)
		return err
	})
	if err != nil {
		panic(err)
	}
}
func evaluate(path string, patterns []string) error {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			panic(err)
		}
		if match {
			fmt.Printf("File '%s' matches pattern '%s'\n",
				path, pattern)
		}
	}
	return nil
}
