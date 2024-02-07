package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type PCStats struct {
	CPULoad        float64 `json:"cpu_load"`
	CPUTemperature float64 `json:"cpu_temperature"`
	TotalMemory    uint64  `json:"total_memory"`
	MemoryUsed     uint64  `json:"memory_used"`
	DiskTotal      uint64  `json:"disk_total"`
	DiskUsage      uint64  `json:"disk_usage"`
	Hostname       string  `json:"hostname"`
	Uptime         uint64  `json:"uptime"`
}

// readCPUTemperature reads the CPU temperature from the first thermal zone found.
func readCPUTemperature() float64 {
	basePath := "/sys/class/thermal/thermal_zone"
	for i := 0; ; i++ {
		path := basePath + strconv.Itoa(i) + "/temp"
		content, err := ioutil.ReadFile(path)
		if err != nil {
			// If we cannot read the temperature, return 0 or handle the error as appropriate.
			return 0
		}
		tempStr := strings.TrimSpace(string(content))
		tempInt, err := strconv.ParseInt(tempStr, 10, 64)
		if err != nil {
			return 0
		}
		// Convert the temperature from millidegree Celsius to degree Celsius.
		return float64(tempInt) / 1000.0
	}
}

func getPCStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	d, _ := disk.Usage("/")
	h, _ := host.Info()

	stats := PCStats{
		CPULoad:        c[0],
		CPUTemperature: readCPUTemperature(), // Updated to use the new function
		TotalMemory:    v.Total,
		MemoryUsed:     v.Used,
		DiskTotal:      d.Total,
		DiskUsage:      d.Used,
		Hostname:       h.Hostname,
		Uptime:         h.Uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	http.HandleFunc("/pcstats", getPCStats)
	http.ListenAndServe("0.0.0.0:9889", nil)
}
