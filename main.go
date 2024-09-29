package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"latte_go_monitor/latte_lib"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Monitor latte_lib.InfluxConfig `json:"monitor_config"`
	LocalIp string                 `json:"local_ip"`
	Pid     int                    `json:"pid"`
	Redis   latte_lib.RedisConfig  `json:"redis_config"`
}

func initLog() {
	file := "./monitor.log"
	err := os.Remove(file)
	if err != nil {
		return
	}
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[log]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func main() {

	//初始化日志  把日志打印到文件内
	initLog()
	configStr := flag.String("config_path", "./monitor.json", "config")
	flag.Parse()
	log.Printf("config: %s", *configStr)
	file, err := os.Open(*configStr)
	if err != nil {
		log.Fatalf("Failed to open file:%v", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read config file content: %v", err)
	}
	file.Close()

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("failed to parse Json: %v", err)
	}
	log.Printf("ip: %s\n", config.LocalIp)
	log.Printf("pid: %d\n", config.Pid)
	log.Printf("redis: %s:%d\n", config.Redis.Host, config.Redis.Port)
	redisC := latte_lib.RedisClient{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		PoolSize: 2,
	}

	redisC.Init()

	influxC := latte_lib.InfluxClient{
		Url: config.Monitor.Url,
		Db:  config.Monitor.Db,
	}
	err = influxC.Init()
	if err != nil {
		log.Printf("influx init fail: %v\n", err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {

		info, err := redisC.GetInfo("cpu", map[string]string{
			"used_cpu_sys":  "float",
			"used_cpu_user": "float",
		})
		if err != nil {
			log.Fatalf("redis get info fail %v", err)
		}

		used_cpu_sys := (info["used_cpu_sys"]).(float64)
		log.Printf("used_cpu_sys: %f\n", used_cpu_sys)
		used_cpu_user := (info["used_cpu_user"]).(float64)
		log.Printf("used_cpu_user: %f\n", used_cpu_user)
		cpus := latte_lib.IoStat(config.Pid)
		log.Printf("sysCpu %f userCpu %f", cpus["total_sys_cpu"], cpus["total_user_cpu"])

		err = influxC.Send("k8s.cache.pid.sys_cpu", map[string]string{
			"pid": strconv.Itoa(config.Pid),
			"ip":  config.LocalIp,
			"idc": "SHA-ALI",
		}, map[string]interface{}{
			"value": cpus["total_sys_cpu"],
		})
		if err != nil {
			log.Printf("influx send pid sys_cpu fail: %v\n", err)
		}

		err = influxC.Send("k8s.cache.pid.user_cpu", map[string]string{
			"pid": strconv.Itoa(config.Pid),
			"ip":  config.LocalIp,
			"idc": "SHA-ALI",
		}, map[string]interface{}{
			"value": cpus["total_user_cpu"],
		})
		if err != nil {
			log.Printf("influx send pid sys_cpu fail: %v\n", err)
		}
	}

}
