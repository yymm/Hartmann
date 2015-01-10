package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Json schema
type json_struct struct {
	Stdout  string
	Stderr  string
	Command string
	Status  int
	Date    string
	App     string
}

// Index
func Index(rw http.ResponseWriter, req *http.Request) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	//get
	tag := req.URL.Path[1:]
	l, err := redis.Int(c.Do("llen", tag))
	if err != nil {
		fmt.Fprintf(rw, "key not found")
		return
	}
	rec, err := redis.Strings(c.Do("lrange", tag, "0", l))
	if err != nil {
		fmt.Fprintf(rw, "key not found")
		return
	}

	for i, _ := range rec {
		fmt.Fprintf(rw, "<div>[%d] json: %s</div>", i, rec[i])
	}
}

// Post: Json data
func PostJson(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t json_struct
	if err := decoder.Decode(&t); err != nil {
		log.Fatal("json.Decode: ", err)
	}
	t.Date = time.Now().Format("2006-01-02T15:04:05Z07:00")
	//log.Printf(t.Stdout + "\n")
	//log.Printf(t.Stderr + "\n")
	//log.Printf("%s => status(%d)\n", t.Command, t.Status)
	SaveJsonToRedis(t.App, t)
	ShowNotifier(t.App, t.Status)
}

func SaveJsonToRedis(app string, j json_struct) {
	data, err := json.Marshal(j)
	if err != nil {
		log.Fatal("json.Marshal: ", err)
	}
	//fmt.Printf(string(data))

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal("redis.Dial: ", err)
	}
	defer c.Close()
	c.Do("rpush", app, string(data))

	log.Printf("[App: %s] success to connect and save to redis.", app)

}

// Call notify tool
// Enable: OSX and Ubuntu
func ShowNotifier(app string, status int) {
	dir, err := filepath.Abs(filepath.Dir("hartmann.jpg"))
	if err != nil {
		log.Fatal(err)
	}

	s := "Success!!"
	if status == 0 {
		s = "Failure.."
	}

	if _, err := os.Stat("/System/Library/CoreServices/SystemVersion.plist"); err == nil {
		out, err := exec.Command("terminal-notifier",
			"-group", "'Hartmann'",
			"-title", "["+s+"]"+"Application: "+app,
			"-subtitle", "Sir! This is a Hartmann notification",
			"-message", s,
			"-contentImage", filepath.Join(dir, "hartmann.jpg"),
			"-appIcon", filepath.Join(dir, "hartmann.jpg")).Output()
		if err != nil {
			log.Fatal("exec.Command: ", err)
		}
		log.Printf("Show notifier: %s", out)
	} else if _, err := os.Stat("/etc/lsb-release"); err == nil {
		out, err := exec.Command("notify-send",
			"["+s+"]"+"Application: "+app, "Sir! This is a Hartmann notification",
			"-i", filepath.Join(dir, "hartmann.jpg")).Output()
		if err != nil {
			log.Fatal("exec.Command: ", err)
		}
		log.Printf("Show notifier: %s", out)
	} else {
		log.Fatal("Invalid os...")
	}
}

func main() {
	port := "8100"
	http.HandleFunc("/", Index)
	http.HandleFunc("/json", PostJson)
	fmt.Printf("Start server. Port: " + port + " (default)\n\n")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
