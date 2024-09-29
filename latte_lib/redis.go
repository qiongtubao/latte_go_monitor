package latte_lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisClient struct {
	Addr     string
	PoolSize int
	client   *redis.Client
}

func (r *RedisClient) Init() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: "",
		DB:       0,
		PoolSize: r.PoolSize,
	})
}

func getStringInfoAttributes(infoContent string, key string) ([]string, error) {
	reg := regexp.MustCompile(`\n` + key + `:(?s:(.*?))\n`)
	if reg == nil {
		return nil, fmt.Errorf("MustCompile info attribute (%s) err", key)
	}
	matchResult := reg.FindAllStringSubmatch(infoContent, -1)
	result := make([]string, len(matchResult))
	for i, v := range matchResult {
		if len(v) < 2 {
			return nil, fmt.Errorf("info attribute match (%s) index (%d) err", key, i)
		}
		result[i] = strings.Trim(v[1], "\r")
	}
	return result, nil

}

func getStringInfoAttribute(info string, key string) (string, error) {
	result, err := getStringInfoAttributes(info, key)
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		return "", errors.New(fmt.Sprintf("info not find (%s) attribute %v", key, result))
	}
	return result[0], nil
}

func getFloatInfoAttribute(info string, key string) (float64, error) {
	strValue, err := getStringInfoAttribute(info, key)
	if err != nil {
		return 0, err
	}
	percentage := 0
	if strValue[len(strValue)-1] == '%' {
		strValue = strValue[:len(strValue)-1]
		percentage = 1
	}
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		log.Printf("info %s=%s parse float fail", key, strValue)
	}
	if percentage == 1 {
		value = value / 100
	}
	return value, err
}

func (r *RedisClient) GetInfo(key string, fields map[string]string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	cpuInfo, err := r.client.Info(ctx, "cpu").Result()
	if err != nil {
		log.Fatalf("redis get info fail %v", err)
		return nil, err
	}
	result := map[string]interface{}{}
	for key := range fields {
		switch fields[key] {
		case "float":
			f, err := getFloatInfoAttribute(cpuInfo, key)
			if err != nil {
				return nil, err
			}
			result[key] = f
		case "string":
			s, err := getStringInfoAttribute(cpuInfo, key)
			if err != nil {
				return nil, err
			}
			result[key] = s
		}
	}
	defer cancel()
	return result, nil
}
