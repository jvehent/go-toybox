#!/usr/bin/env python
# ulfr - 2013

import hashlib
from re import compile, match
from collections import defaultdict
from optparse import OptionParser
from time import sleep
import os
from sys import exit

# FIXME: don't pass global dict around, use 'return' instead
CHECKLIST = {}
""" CHECKLIST contains a dict of files with associated checks """
RESULTS = {}
""" RESULTS contains a dict of IOCs with positive and negative matches """
counter = 0
""" Keep trac of things """

def build_dict_of_checks(checks, mode):
	d = {}
	for ioc in checks:
		if checks[ioc]['mode'] == mode:
			d[ioc] = checks[ioc]['check']
	return d


def compute_md5_of_file(f):
	"""
		Return the MD5 hash of a file, computed block by block,
	"""
	BLOCK_SIZE = 128
	md5 = hashlib.md5()
	if options.debug:
		print("[DEBUG] compute_md5_of_file: called with file '%s'" % (f,))
	try:
		fd = open(f)
	except IOError:
		if options.verbose:
			print("[ERROR] file_has_md5: Can't open '%s'" % (f,))
		return None
	while True:
		data = fd.read(BLOCK_SIZE)
		if not data:
			break
		md5.update(data)
	fd.close()
	return md5.hexdigest()


def add_file_to_checklist(f, ioc, mode, check):
	"""
		Add a file to the global checklist
	"""
	if options.debug:
		print("Adding file '%s' mode '%s' check '%s'" % (f, mode, check))
	if f not in CHECKLIST:
		CHECKLIST[f] = {}
	modes = CHECKLIST[f].get('modes', [])
	if mode not in modes:
		modes.append(mode)
		CHECKLIST[f]['modes'] = modes
	iocs = CHECKLIST[f].get('iocs', {})
	if mode == 'regex':
		""" compile the regex	first, then store it """
		regex = compile(check)
		iocs[ioc] = { 'mode': mode, 'check': regex }
	else:
		iocs[ioc] = { 'mode': mode, 'check': check }
	CHECKLIST[f]['iocs'] = iocs


def store_match(f, ioc, counter):
	CHECKLIST[f]['iocs'][ioc]['result'] = True
	if ioc not in RESULTS:
		RESULTS[ioc] = []
	RESULTS[ioc].append(f)
	if options.debug:
		print("[DEBUG] IOC '%s' found match on file '%s'" % (ioc, f))
	return counter + 1


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

for ioc in args:
	"""
		IOCs are stored in positional arguments, format is <path>:<mode>=<check>
		This loop parses the arguments and build the list of files to check
	"""
	path = ioc.split(':')[0]
	(mode, check) = ioc.split(':')[1].split('=')
	if mode not in ['regex', 'md5', 'sha1']:
		print("[ERROR] Unknown IOC mode '%s'" % (mode,))
		sys.exit(1)
	if os.path.isdir(path):
		for r,d,flist in os.walk(path):
			for f in flist:
				f_abs = os.path.join(r,f)
				add_file_to_checklist(f_abs, ioc, mode, check)
	elif os.path.isfile(path):
		add_file_to_checklist(path, ioc, mode, check)

for f in CHECKLIST:
	"""
		Check the files listed in the CHECKLIST dict.
	"""
	if 'md5' in CHECKLIST[f]['modes']:
		h = compute_md5_of_file(f)
		md5s = build_dict_of_checks(CHECKLIST[f]['iocs'], 'md5')
		for ioc in md5s:
			if h == md5s[ioc]:
				counter = store_match(f, ioc, counter)
	if 'regex' in CHECKLIST[f]['modes']:
		regexes = build_dict_of_checks(CHECKLIST[f]['iocs'], 'regex')
		try:
			fd = open(f)
		except IOError:
			if options.verbose:
				print("[ERROR] Can't open file '%s'" % f)
		for line in fd:
			for ioc in regexes:
				if regexes[ioc].match(line):
					counter = store_match(f, ioc, counter)

print("%s results found." % (counter,))
for ioc in RESULTS:
	print("IOC '%s' present in '%s'" % (ioc, RESULTS[ioc]))
