#!/usr/bin/env python

import urllib2
import json
import subprocess
import sys


def exec_command(cmdline):
    p = subprocess.Popen(cmdline, shell=True,
                         cwd='.',
                         stdin=subprocess.PIPE,
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE,
                         close_fds=True)
    return p.stdout, p.stdin, p.stderr


if len(sys.argv) != 3:
    print("Argvs is two.")
    exit(1)

app = sys.argv[1]
cmd = sys.argv[2]
stdout, stdin, stderr = exec_command(cmd)

stdout_utf8 = stdout.read().decode('utf-8')
stderr_utf8 = stderr.read().decode('utf-8')

json_dic = {
    'stdout': stdout_utf8,
    'stderr': stderr_utf8,
    'status': 0 if len(stderr_utf8) else 1,
    'command': cmd,
    'app': app
}

json_str = json.dumps(json_dic)

#print(json_str)

headers = {'Content-type': 'application/json'}

req = urllib2.Request('http://192.168.5.52:8100/json', json_str, headers)

f = urllib2.urlopen(req)

response = f.read()
f.close()
