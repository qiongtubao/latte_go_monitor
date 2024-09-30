package latte_lib

/*
type TimesStat struct {
	User        float64 `json:"user"`       // 用户态消耗的时间
	Nice        float64 `json:"nice"`       // 用户态下优先级提升的消耗时间
	System      float64 `json:"system"`     // 内核态消耗的时间
 Idle        float64 `json:"idle"`        // CPU 空闲的时间
	Iowait      float64 `json:"iowait"`     // 等待 I/O 完成的时间
	Irq         float64 `json:"irq"`        // 处理硬件中断的时间
	Softirq     float64 `json:"softirq"`    // 处理软件中断的时间
	Steal       float64 `json:"steal"`      // 被虚拟机偷走的时间
	Guest       float64 `json:"guest"`      // 在用户模式下运行的虚拟 CPU 时间
	GuestNice   float64 `json:"guest_nice"` // 在用户模式下运行的虚拟 CPU 时间（优先级提升）
}
*/

// func calculatePercent(t1, t2 float64, delta float64, numcpu int) float64 {
// 	if delta == 0 {
// 		return 0
// 	}
// 	delta_proc := t2 - t1
// 	overall_percent := ((delta_proc / delta) * 100) * float64(numcpu)
// 	return overall_percent
// }

// func CpuStat(pid int, wait int) map[string]float64 {
// 	// 获取进程对象
// 	p, err := process.NewProcess(int32(pid))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	start_now := time.Now()
// 	numcpu := runtime.NumCPU()
// 	// 获取第一次 CPU 时间
// 	start_times, err := p.Times()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// 等待一段时间
// 	// sysUsage, _ := p.Percent(1 * time.Second)
// 	time.Sleep(time.Duration(wait) * time.Second)

// 	// 获取第二次 CPU 时间
// 	end_times, err := p.Times()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	now := time.Now()
// 	delta := (now.Sub(start_now).Seconds()) * float64(numcpu)
// 	// 计算用户态和内核态的 CPU 使用率
// 	userCpuUsage := calculatePercent(start_times.User, end_times.User, delta, numcpu)
// 	systemCpuUsage := calculatePercent(start_times.System, end_times.System, delta, numcpu)
// 	totalCpuUsage := calculatePercent(start_times.Total(), end_times.Total(), delta, numcpu)
// 	// log.Printf("user usage %f sys usage %f  psys usage %f", userCpuUsage, systemCpuUsage, sysUsage)
// 	return map[string]float64{
// 		"sys_cpu":   systemCpuUsage,
// 		"user_cpu":  userCpuUsage,
// 		"total_cpu": totalCpuUsage,
// 	}
// }
