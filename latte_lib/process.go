package latte_lib

import (
	"log"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type LatteProcess struct {
	Pid int
	p   *process.Process
}

func (lp *LatteProcess) Init() error {
	p, err := process.NewProcess(int32(lp.Pid))
	if err != nil {
		log.Fatal(err)
		return err
	}
	lp.p = p
	return nil
}
func calculate(s float64, e float64, delta float64) float64 {
	return (e - s) / delta
}
func (lp *LatteProcess) GetIOStats(time_interval time.Duration) map[string]float64 {
	// numcpu := runtime.NumCPU()
	start_time := time.Now()
	// 获取第一次 磁盘I/O 信息
	start_io_stats, err := lp.p.IOCounters()
	if err != nil {
		log.Fatal(err)
	}
	// 等待一段时间
	time.Sleep(time_interval)

	// 获取第二次 磁盘I/O 信息
	end_io_stats, err := lp.p.IOCounters()
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	delta := (now.Sub(start_time).Seconds())
	read_iops := calculate(float64(start_io_stats.ReadCount), float64(end_io_stats.ReadCount), delta)
	write_iops := calculate(float64(start_io_stats.WriteCount), float64(end_io_stats.WriteCount), delta)
	read_throughput := calculate(float64(start_io_stats.ReadBytes), float64(end_io_stats.ReadBytes), delta)
	write_throughput := calculate(float64(start_io_stats.WriteBytes), float64(start_io_stats.WriteBytes), delta)

	return map[string]float64{
		"read_iops":        read_iops,
		"write_iops":       write_iops,
		"read_throughput":  read_throughput,
		"write_throughput": write_throughput,
	}
}

func calculateCpuPercent(t1, t2 float64, delta float64, numcpu int) float64 {
	if delta == 0 {
		return 0
	}
	delta_proc := t2 - t1
	overall_percent := ((delta_proc / delta) * 100) * float64(numcpu)
	return overall_percent
}

func (lp *LatteProcess) GetCpuStats(time_interval time.Duration) map[string]float64 {
	numcpu := runtime.NumCPU()
	start_time := time.Now()
	// 获取第一次 CPU 时间
	start_times, err := lp.p.Times()
	if err != nil {
		log.Fatal(err)
	}

	// 等待一段时间
	time.Sleep(time_interval)

	// 获取第二次 CPU 时间
	end_times, err := lp.p.Times()
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	delta := (now.Sub(start_time).Seconds()) * float64(numcpu)
	// 计算用户态和内核态的 CPU 使用率
	userCpuUsage := calculateCpuPercent(start_times.User, end_times.User, delta, numcpu)
	systemCpuUsage := calculateCpuPercent(start_times.System, end_times.System, delta, numcpu)
	totalCpuUsage := calculateCpuPercent(start_times.Total(), end_times.Total(), delta, numcpu)
	// log.Printf("user usage %f sys usage %f  psys usage %f", userCpuUsage, systemCpuUsage, sysUsage)
	return map[string]float64{
		"sys_cpu":   systemCpuUsage,
		"user_cpu":  userCpuUsage,
		"total_cpu": totalCpuUsage,
	}
}
