package server

import (
	"fmt"
	"log"
	"net/http"

	. "common"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var (
	emptyVal               = struct{}{}
	IGNORE_DISK_FSTYPE_MAP = map[string]struct{}{
		"devtmpfs":    emptyVal,
		"autofs":      emptyVal,
		"sysfs":       emptyVal,
		"proc":        emptyVal,
		"tmpfs":       emptyVal,
		"binfmt_misc": emptyVal,
		"cgroupfs":    emptyVal,
		"hugetlbfs":   emptyVal,
		"mqueue":      emptyVal,
		"fusectl":     emptyVal,
		"devpts":      emptyVal,
		"securityfs":  emptyVal,
		"cgroup":      emptyVal,
		"pstore":      emptyVal,
		"debugfs":     emptyVal,
		"selinuxfs":   emptyVal,
	}
)

func SystemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !CheckAuth(vars["key"], w, r) {
		return
	}
	memState, memOut := checkMem(vars["mem"])
	diskState, diskOut := checkDisk(vars["disk"])
	fmt.Fprintf(w, "%d|%s[%s], %s[%s]", HighState(memState, diskState), memOut, StateString(memState), diskOut, StateString(diskState))
}

func checkDisk(param string) (state int, out string) {
	out, state = "", OK
	diskW, diskC, err := NagiosParam(param)
	parts, err := disk.Partitions(false)
	if err != nil {
		log.Println(err)
		out = "DISK:ERR"
		state = UNKNOWN
		return
	}
	out = "DISK:"
	for _, part := range parts {
		if _, ok := IGNORE_DISK_FSTYPE_MAP[part.Fstype]; ok {
			continue
		}
		v, err := disk.Usage(part.Mountpoint)
		if err != nil {
			log.Println(err)
			out = "DISK:ERR"
			state = UNKNOWN
			return
		}
		usedPercent := (1.0 - float64(v.Free)/float64(v.Total)) * 100
		if usedPercent > diskC {
			state = CRITICAL
		} else if usedPercent > diskW {
			if state != CRITICAL {
				state = WARNING
			}
		}
		out = fmt.Sprintf("%s[%s-%s]%.1f,", out, part.Mountpoint, part.Fstype, usedPercent)
	}
	return
}

func checkMem(param string) (state int, out string) {
	out, state = "", OK
	memW, memC, err := NagiosParam(param)
	if err != nil {
		log.Println(err)
		out = "MEM:INVP"
		state = UNKNOWN
		return
	}

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println("failure to fetch memory", err)
		out = "MEM:ERR"
		state = UNKNOWN
		return
	}
	out = fmt.Sprintf("MEM:%.1f", v.UsedPercent)
	if v.UsedPercent > memC {
		state = CRITICAL
	} else if v.UsedPercent > memW {
		state = WARNING
	}
	return
}
