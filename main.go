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

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/pkg/errors"
)

type Config struct {
	Monitor       latte_lib.InfluxConfig `json:"monitor_config"`
	LocalIp       string                 `json:"local_ip"`
	Pid           int                    `json:"pid"`
	Time_interval int                    `json:"time_interval"`
	Redis         latte_lib.RedisConfig  `json:"redis_config"`
}

func setupLogRotation(file string, rotationTime time.Duration, maxSize int64) (writer *rotatelogs.RotateLogs, err error) {
	hook, err := rotatelogs.New(
		file+".%Y%m%d",
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithLinkName(file),        // symlink current log file
		rotatelogs.WithRotationCount(10),     // keep 5 old log files
		rotatelogs.WithRotationSize(maxSize), // max size of each log file in bytes
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rotatelogs instance")
	}

	writer = hook
	return
}
func initLog() {
	file := "./monitor.log"
	writer, err := setupLogRotation(file, 48*time.Hour, 100*1024*1024)
	if err != nil {
		log.Fatalf("failed to set up log rotation:%v", err)
		return
	}
	log.SetOutput(writer) // 将文件设置为log输出的文件
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

	ticker := time.NewTicker(time.Duration(config.Time_interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {

		// info, err := redisC.GetInfo("cpu", map[string]string{
		// 	"used_cpu_sys":  "float",
		// 	"used_cpu_user": "float",
		// })
		// if err != nil {
		// 	log.Fatalf("redis get info fail %v", err)
		// }

		// used_cpu_sys := (info["used_cpu_sys"]).(float64)
		// log.Printf("used_cpu_sys: %f\n", used_cpu_sys)
		// used_cpu_user := (info["used_cpu_user"]).(float64)
		// log.Printf("used_cpu_user: %f\n", used_cpu_user)
		cpus := latte_lib.IoStat(config.Pid, config.Time_interval-1)
		log.Printf("sysCpu %f userCpu %f totalCpu %f", cpus["sys_cpu"], cpus["user_cpu"], cpus["total_cpu"])

		err = influxC.Send("k8s.cache.pid.sys_cpu", map[string]string{
			"pid": strconv.Itoa(config.Pid),
			"ip":  config.LocalIp,
			"idc": "SHA-ALI",
		}, map[string]interface{}{
			"value": cpus["sys_cpu"],
		})
		if err != nil {
			log.Printf("influx send pid sys_cpu fail: %v\n", err)
		}

		err = influxC.Send("k8s.cache.pid.user_cpu", map[string]string{
			"pid": strconv.Itoa(config.Pid),
			"ip":  config.LocalIp,
			"idc": "SHA-ALI",
		}, map[string]interface{}{
			"value": cpus["user_cpu"],
		})
		if err != nil {
			log.Printf("influx send pid sys_cpu fail: %v\n", err)
		}

		err = influxC.Send("k8s.cache.pid.total_cpu", map[string]string{
			"pid": strconv.Itoa(config.Pid),
			"ip":  config.LocalIp,
			"idc": "SHA-ALI",
		}, map[string]interface{}{
			"value": cpus["total_cpu"],
		})
		if err != nil {
			log.Printf("influx send pid sys_cpu fail: %v\n", err)
		}
	}

}
