package client

import (
	. "common"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	key                string
	host               string
	port               int
	system             bool
	systemDiskWarning  int
	systemDiskCritical int
	systemMemWarning   int
	systemMemCritical  int
	SYSTEM_API         = "http://%s:%d/%s/system/disk/%d|%d/mem/%d|%d"
)

func init() {
	flag.StringVar(&key, "k", "test12345", "nagios agent key")
	flag.StringVar(&host, "h", "127.0.0.1", "nagios agent host")
	flag.IntVar(&port, "p", 16888, "nagios agent port")
	flag.BoolVar(&system, "system", false, "monitor system status")
	flag.IntVar(&systemDiskWarning, "systemDiskWarning", 80, "system disk warning value")
	flag.IntVar(&systemDiskCritical, "systemDiskCritical", 90, "system disk critical value")
	flag.IntVar(&systemMemWarning, "systemMemWarning", 80, "system memory warning value")
	flag.IntVar(&systemMemCritical, "systemMemCritical", 90, "system memory critical value")
	flag.Parse()
}
func Run() {
	if system {
		systemApi := fmt.Sprintf(SYSTEM_API, host, port, key, systemDiskWarning, systemDiskCritical, systemMemWarning, systemMemCritical)
		state, msg := InvokeApi(systemApi)
		fmt.Println(msg)
		os.Exit(state)
	}
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(UNKNOWN)
}

func InvokeApi(url string) (state int, msg string) {
	var err error
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return UNKNOWN, err.Error()
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return UNKNOWN, err.Error()
	}
	content := string(body)
	if len(content) > 3 {
		if state, err = strconv.Atoi(content[:1]); err != nil {
			return UNKNOWN, err.Error()
		}
		sep := content[1:2]
		if sep != "|" {
			return UNKNOWN, fmt.Sprintf("invalid api content: %s", content)
		}
		msg = content[2:]
		return
	}
	return UNKNOWN, fmt.Sprintf("invalid api content: %s", content)
}
