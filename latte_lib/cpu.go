package latte_lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func IoStat(pid int) map[string]float64 {
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	content, err := ioutil.ReadFile(statFile)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", statFile, err)
	}
	//解析 stat 文件内容
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	scanner.Scan()
	statLine := scanner.Text()

	fields := strings.Fields(statLine)
	userTime, _ := strconv.ParseFloat(fields[13], 64)
	systemTime, _ := strconv.ParseFloat(fields[14], 64)
	childrenUserTime, _ := strconv.ParseFloat(fields[15], 64)
	childrenSystemTime, _ := strconv.ParseFloat(fields[16], 64)
	totalSysTime := systemTime + childrenSystemTime
	totalUserTime := userTime + childrenUserTime
	// totalTime := userTime + systemTime + childrenUserTime + childrenSystemTime

	sysTotalFile := "/proc/stat"
	sysContent, err := ioutil.ReadFile(sysTotalFile)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", sysTotalFile, err)
	}

	sysFields := strings.Fields(string(sysContent))
	sysTotal, _ := strconv.ParseFloat(sysFields[1], 64)

	return map[string]float64{
		"sys_cpu":        systemTime / 100,
		"user_cpu":       userTime / 100,
		"total_sys_cpu":  totalSysTime / 100,
		"total_user_cpu": totalUserTime / 100,
		"total_cpu":      sysTotal,
	}
}
