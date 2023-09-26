// Package main
// @Author Joe
// @Date 2023-09-23 20:47:33
// @Description:
package connection

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"redis-dump/utils"
	"strconv"
	"strings"
	"time"
)

var RedisPool redis.Pool

/*RedisParams
 * @Description:
 * @Author Joe
 * @Date 2023-09-23 16:05:09
 */
type RedisParams struct {
	Db    int         `json:"db"`
	Key   string      `json:"key"`
	TTL   int         `json:"ttl"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	Size  int         `json:"size"`
}

/**
 * initRedisPool
 * @Description 初始化reids连接池
 * @Author Joe
 * @Date 2023-09-22 10:50:17
 * @Param ip string
 * @Param port int
 * @Param password string
 * @Param dbNum int
 * @return conn redis.Pool
 **/
func InitRedisPool(ip string, port int, password string, dbNum int) (conn redis.Pool) {
	return redis.Pool{
		MaxIdle:     0,                  // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive:   0,                  // 最大的连接数，表示同时最多有N个连接。0表示不限制。
		IdleTimeout: time.Duration(120), // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭。如果设置成0，空闲连接将不会被关闭。应该设置一个比redis服务端超时时间更短的时间。
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				ip+":"+strconv.Itoa(port),
				redis.DialReadTimeout(time.Duration(10000)*time.Millisecond),    // 从Redis读取数据超时时间。
				redis.DialWriteTimeout(time.Duration(10000)*time.Millisecond),   // 向Redis写入数据超时时间。
				redis.DialConnectTimeout(time.Duration(10000)*time.Millisecond), // 连接Redis超时时间。
				redis.DialDatabase(dbNum),    // select 库
				redis.DialPassword(password), // 密码
			)
		},
	}
}

/**
 * getRedisDbInfo
 * @Description 获取redis Db信息
 * @Author Joe
 * @Date 2023-09-22 10:49:37
 * @return map[int]int
 * @return error
 **/
func GetRedisDbInfo() (map[int]int, error) {
	result := GetRedisInfo()

	if _, ok := result["Keyspace"]; !ok {
		return nil, errors.New("未找到 Keyspace")
	}

	if len(result["Keyspace"]) <= 0 {
		return nil, errors.New("Keyspace 为空")
	}

	dbInfo := make(map[int]int)
	for key, val := range result["Keyspace"] {

		// 以 "," 分隔键和值
		parts := strings.Split(val, ",")
		keyInfo := strings.TrimSpace(parts[0])
		// 以 "=" 分隔键和值
		parts = strings.Split(keyInfo, "=")
		keyCount := strings.TrimSpace(parts[1])

		dbNum, _ := strconv.Atoi(key[2:])
		keyCountInt, _ := strconv.Atoi(keyCount)
		dbInfo[dbNum] = keyCountInt
	}

	return dbInfo, nil
}

/**
 * getRedisInfo
 * @Description 获取redis信息
 * @Author Joe
 * @Date 2023-09-22 10:50:30
 * @return map[string]map[string]string
 **/
func GetRedisInfo() map[string]map[string]string {
	conn := RedisPool.Get()
	defer conn.Close()
	// 获取Redis的信息
	info, err := redis.String(conn.Do("INFO"))
	if err != nil {
		utils.WriteLog("", err.Error())
	}

	result := make(map[string]map[string]string)
	currentKey := ""

	lines := strings.Split(info, "\n")
	infoKey := ""
	for _, line := range lines {
		// 以 "#" 开头的行被认为是注释，跳过处理
		if strings.HasPrefix(line, "#") {
			infoKey = line[2 : len(line)-1]
			result[infoKey] = map[string]string{}
			continue
		}

		// 以 ":" 分隔键和值
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 如果新的键进行了指定，则更新 currentKey
		if key != "" {
			currentKey = key
		}

		result[infoKey][currentKey] = value
	}

	return result
}
