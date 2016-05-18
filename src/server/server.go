package server

import (
	. "common"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

const (
	MAX_TIMES_TRY_PER_IP = 50
)

var (
	key      string
	addr     string
	demonize bool
	ipMap    = make(map[string]int)
)

func init() {
	flag.BoolVar(&demonize, "d", false, "demonize")
	flag.StringVar(&key, "k", "test12345", "auth key")
	flag.StringVar(&addr, "addr", ":16888", "server listen address")
	flag.Parse()
}

func Run() {
	if key == "" || addr == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	if demonize {
		go handleSignal()
	}
	r := mux.NewRouter()
	r.HandleFunc("/{key}/system/disk/{disk}/mem/{mem}", SystemHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for sig := range c {
		switch sig {
		case syscall.SIGHUP:
		}
	}
}
func CheckAuth(inKey string, w http.ResponseWriter, r *http.Request) bool {
	ip := r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")]
	count, ok := ipMap[ip]
	if ok && count > MAX_TIMES_TRY_PER_IP {
		fmt.Fprintf(w, "%d|max retry time reached:%s", UNKNOWN, inKey)
		return false
	}
	if inKey != key {
		ipMap[ip] = count + 1
		fmt.Fprintf(w, "%d|invalid key:%s", UNKNOWN, inKey)
		return false
	}
	return true
}
