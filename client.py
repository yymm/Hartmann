import http.client
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


connection = http.client.HTTPConnection('localhost', 8100)

headers = {'Content-type': 'application/json'}

if len(sys.argv) != 3:
    print("Argvs is two.")
    exit(1)

app = sys.argv[1]
cmd = sys.argv[2]
stdout, stdin, stderr = exec_command(cmd)

stdout_utf8 = stdout.read().decode('utf-8')
stderr_utf8 = stderr.read().decode('utf-8')

#print(stdout_utf8)
#print(stderr_utf8)
#print(cmd)

json_dic = {
    'stdout': stdout_utf8,
    'stderr': stderr_utf8,
    'status': 0 if len(stderr_utf8) else 1,
    'command': cmd,
    'app': app
}

print(json_dic)

json_str = json.dumps(json_dic)

print(json_str)

connection.request('POST', '/json', json_str, headers)

response = connection.getresponse()
print(response.read().decode())
