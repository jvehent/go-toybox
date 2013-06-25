#!/usr/bin/env python

import hashlib
from re import compile, match
from optparse import OptionParser
from time import sleep
import os


def file_has_md5(f, h):
	"""
		Calculate the MD5 hash of a file, block by block,
		and return True if it matches the submitted hash
	"""
	BLOCK_SIZE = 128
	md5 = hashlib.md5()
	if options.debug:
		print("[DEBUG] file_has_md5: hash '%s' file '%s'" % (h, f))
	try:
		fd = open(f)
	except IOError:
		if options.verbose:
			print("[ERROR] file_has_md5: Can't open '%s'" % f)
	while True:
		data = fd.read(BLOCK_SIZE)
		if not data:
			break
		md5.update(data)
	fd.close()
	if h == md5.hexdigest():
		if options.debug:
			print("[DEBUG] file_has_md5: match found file '%s' hash '%s'" %
				  (f, h))
		return True
	return False


def file_has_regex(f, regex):
	"""
		Look for a matching regex inside a file, line by line,
		and return True if found
	"""
	fd = open(f)
	if options.debug:
		print("[DEBUG] file_has_regex: checking file '%s'" % (f,))
	try:
		fd = open(f)
	except IOError:
		if options.verbose:
			print("[ERROR] file_has_regex: Can't open '%s'" % f)
	for line in fd:
		if regex.match(line):
			if options.debug:
				print("[DEBUG] file_has_regex: match found in file '%s'" %
					  (f,))
			return True
	return False


usage = """%prog [options] <path>:regex=<regex> <path>:md5=<hash> ...
Compare the files contained in a path with a regex or a hash."""
parser = OptionParser(usage=usage)
parser.disable_interspersed_args()
parser.add_option("-o", "--output", dest="report", type="string",
                  help="write results to REPORT file")
parser.add_option("-t", "--throttle", dest="throttle", type="float", default=0,
                  help="Sleep for X seconds between each files (ex: 0.01)")
parser.add_option("-v", "--verbose",
                  action="store_true", dest="verbose", default=True,
                  help="make lots of noise [default]")
parser.add_option("-q", "--quiet",
                  action="store_false", dest="verbose",
                  help="be vewwy quiet (I'm hunting wabbits)")
parser.add_option("-d", "--debug",
                  action="store_true", dest="debug", default=False,
                  help="debug logging, to impress yours friends")
(options, args) = parser.parse_args()

if options.verbose:
	print "Starting search. Target list: %s" % args

IOCS = {}
for ioc in args:
	""" IOCs are stored in positional arguments
		the format is <file path>:<mode>=<check value>
	"""
	path = ioc.split(':')[0]
	mode = ioc.split(':')[1].split('=')[0]
	check = ioc.split(':')[1].split('=')[1]
	IOCS[ioc] = {'path': path, 'mode': mode, 'check': check}
	if options.verbose:
		print("Adding IOC '%s' with mode '%s' check '%s'" %
			  (IOCS[ioc]['path'], mode, check))

for ioc in IOCS:
	"""
		Build list of files from the path in the IOC
	"""
	path = IOCS[ioc]['path']
	IOCS[ioc]['files'] = {}
	if os.path.isdir(path):
		for r,d,flist in os.walk(path):
			for f in flist:
				IOCS[ioc]['files'][os.path.join(r,f)] = None
				if options.debug:
					print("Adding file '%s' to target list for IOC '%s'" %
						  (os.path.join(r,f), ioc))
	elif os.path.isfile(path):
		IOCS[ioc]['files'][path] = None

for ioc in IOCS:
	"""
		For each file listed in each IOC, perform the check
		according to the mode
		'regex' mode executes the regex inside the file
		'md5' and 'sha1' modes calculate a hash of the file and compare it to
		the provided hash value
	"""
	if options.verbose:
		print("Checking IOC '%s' against %s files" %
			  (ioc, len(IOCS[ioc]['files'])))
	mode = IOCS[ioc]['mode']
	check = IOCS[ioc]['check']
	for f in IOCS[ioc]['files']:
		if options.throttle > 0:
			sleep(options.throttle)
		if options.debug:
			print("[DEBUG] checking '%s' against '%s' in mode '%s'" %
				  (f, check, mode))
		if mode == 'regex':
			regex = compile(check)
			IOCS[ioc]['files'][f] = file_has_regex(f, regex)
			continue
		if mode == 'md5':
			IOCS[ioc]['files'][f] = file_has_md5(f, check)
			continue

for ioc in IOCS:
	for f in IOCS[ioc]['files']:
		if IOCS[ioc]['files'][f]:
			if options.verbose:
				print("POSITIVE MATCH FOUND: '%s' matches '%s'" %
					  (f, ioc))
