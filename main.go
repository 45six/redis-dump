/**
 * @auther: Joe
 * @date: 2023/9/18 5:07 下午
 *
 */
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"redis-dump/connection"
	"redis-dump/utils"
	"sync"
	"time"
)

// 定义一个同步等待的组
var wg sync.WaitGroup

var operation string
var goroutineChan = make(chan struct{})
var mutex sync.Mutex
var paramMaxCount = 10000

type DataInfo struct {
	DbToKey            map[int]int
	DbToKeyTypeToCount map[int]map[string]int
}

var dataInfo DataInfo

func main() {

	startTime := time.Now()

	helpFlag := flag.Bool("help", false, "Show help\n显示帮助文档")
	ipFlag := flag.String("ip", "127.0.0.1", "redis ip\nredis ip地址")
	protFlag := flag.Int("port", 6379, "redis port\nredis 端口")
	authFlag := flag.String("auth", "", "redis auth\nredis 密码")
	fileFlag := flag.String("file", "redis_data.json", "file path\n文件路径")
	goroutineFlag := flag.Int("gcount", 100, "goroutine count by redis key\n并发执行key操作的最大协程数量")
	operationFlag := flag.String("op", "", "operation(dump or update)\n操作类型(dump备份)(update更新)")
	debugFlag := flag.Bool("debug", false, "is debug\ndebug开关")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("General redis-dump options:")
		flag.PrintDefaults()
		utils.CheckExit()
		return
	}

	flag.Parse()

	if *helpFlag {
		fmt.Println("General redis-dump options:")
		flag.PrintDefaults()
		utils.CheckExit()
		return
	}

	if *operationFlag == "" {
		fmt.Println("missing parameter -op\n")
		utils.CheckExit()
		return
	}

	utils.Debug = *debugFlag
	operation = *operationFlag
	utils.DataFile = *fileFlag
	goroutineChan = make(chan struct{}, *goroutineFlag)

	connection.RedisPool = connection.InitRedisPool(*ipFlag, *protFlag, *authFlag, 0)

	conn := connection.RedisPool.Get()
	_, err := conn.Do("ping")
	if err != nil {
		utils.WriteLog("", err.Error())
		utils.CheckExit()
		return
	}
	conn.Close()

	dataInfo = DataInfo{
		DbToKey:            make(map[int]int),
		DbToKeyTypeToCount: make(map[int]map[string]int),
	}
	switch operation {
	case "dump":
		dump()
	case "update":
		update()
	default:
		fmt.Println("undefined operation: " + operation + ", please enter dump or update")
	}

	elapsed := time.Since(startTime)
	fmt.Println("expend time: ", elapsed)
	utils.CheckExit()
}

/**
 * dump
 * @Description 转储
 * @Author Joe
 * @Date 2023-09-22 10:49:55
 **/
func dump() {
	dbInfo, err := connection.GetRedisDbInfo()
	if err != nil {
		utils.WriteLog("", err.Error())
		return
	}

	conn := connection.RedisPool.Get()
	defer conn.Close()

	for dbNum, keyCount := range dbInfo {
		if keyCount <= 0 {
			continue
		}

		mutex.Lock()
		dataInfo.DbToKey[dbNum] = 0
		dataInfo.DbToKeyTypeToCount[dbNum] = map[string]int{}
		mutex.Unlock()

		conn.Do("SELECT", dbNum)
		keys, err := redis.Strings(conn.Do("KEYS", "*"))
		if err != nil {
			utils.WriteLog("", err.Error())
			return
		}

		for _, redisKey := range keys {
			wg.Add(1)
			goroutineChan <- struct{}{}
			go dumpRedis(dbNum, redisKey)
		}
	}
	//fmt.Println(redisJson)
	//阻塞直到所有任务完成
	wg.Wait()
	for key, val := range dataInfo.DbToKey {
		//fmt.Printf("db number:%v\tkey count:%v\n", key, val)
		utils.WriteLog("", fmt.Sprintf("db number:%v\tkey count:%v", key, val))
		for k, v := range dataInfo.DbToKeyTypeToCount[key] {
			utils.WriteLog("", fmt.Sprintf("key type:%v\tkey count:%v", k, v))
		}
		fmt.Println()
	}
}

func dumpRedis(dbNum int, redisKey string) {
	conn := connection.RedisPool.Get()
	defer func() {
		<-goroutineChan
	}()
	defer wg.Done()
	defer conn.Close()
	conn.Do("SELECT", dbNum)
	keyType, err := redis.String(conn.Do("TYPE", redisKey))
	if err != nil {
		utils.WriteLog("", err.Error())
		return
	}

	keyTtl, err := redis.Int(conn.Do("TTL", redisKey))
	if err != nil {
		utils.WriteLog("", err.Error())
		return
	}

	var keyData connection.RedisParams
	keySize := 0
	switch keyType {
	case "hash":
		stringMap, err := redis.StringMap(conn.Do("HGETALL", redisKey))
		if err != nil {
			utils.WriteLog("", err.Error(), redisKey)
			return
		}

		for key, val := range stringMap {
			keySize += len(key) + len(val)
		}

		keyData = connection.RedisParams{
			Db:    dbNum,
			Key:   redisKey,
			TTL:   keyTtl,
			Type:  keyType,
			Value: stringMap,
			Size:  keySize,
		}

	case "string":
		str, err := redis.String(conn.Do("GET", redisKey))
		if err != nil {
			utils.WriteLog("", err.Error())
			return
		}
		keySize += len(str)
		keyData = connection.RedisParams{
			Db:    dbNum,
			Key:   redisKey,
			TTL:   keyTtl,
			Type:  keyType,
			Value: str,
			Size:  keySize,
		}
	case "set":
		list, err := redis.Strings(conn.Do("SMEMBERS", redisKey))
		if err != nil {
			utils.WriteLog("", err.Error())
			return
		}
		for _, val := range list {
			keySize += len(val)
		}

		keyData = connection.RedisParams{
			Db:    dbNum,
			Key:   redisKey,
			TTL:   keyTtl,
			Type:  keyType,
			Value: list,
			Size:  keySize,
		}
	case "zset":
		stringMap, err := redis.StringMap(conn.Do("ZRANGE", redisKey, 0, -1, "WITHSCORES"))
		if err != nil {
			utils.WriteLog("", err.Error())
			return
		}

		list := [][]string{}
		for key, val := range stringMap {
			list = append(list, []string{key, val})
			keySize += len(key) + len(val)
		}

		keyData = connection.RedisParams{
			Db:    dbNum,
			Key:   redisKey,
			TTL:   keyTtl,
			Type:  keyType,
			Value: list,
			Size:  keySize,
		}
	case "list":
		list, err := redis.Strings(conn.Do("LRANGE", redisKey, 0, -1))
		if err != nil {
			utils.WriteLog("", err.Error())
			return
		}

		for _, val := range list {
			keySize += len(val)
		}

		keyData = connection.RedisParams{
			Db:    dbNum,
			Key:   redisKey,
			TTL:   keyTtl,
			Type:  keyType,
			Value: list,
			Size:  keySize,
		}
	default:
		utils.WriteLog("", fmt.Sprintf("未知类型：%v | %v ", keyType, redisKey))
		return
	}

	keyDataJson, _ := json.Marshal(keyData)
	utils.WriteData("", utils.DataFile, string(keyDataJson)+"\n")

	mutex.Lock()
	dataInfo.DbToKey[dbNum] += 1
	if _, ok := dataInfo.DbToKeyTypeToCount[dbNum][keyType]; !ok {
		dataInfo.DbToKeyTypeToCount[dbNum][keyType] = 1
	} else {
		dataInfo.DbToKeyTypeToCount[dbNum][keyType] += 1
	}
	mutex.Unlock()
}

/**
 * update
 * @Description 更新
 * @Author Joe
 * @Date 2023-09-22 10:50:00
 **/
func update() {

	// 打开文件
	file, err := os.Open(utils.DataFile)
	if err != nil {
		utils.WriteLog("", err.Error())
		return
	}
	defer file.Close()

	// 创建一个新的 Reader 对象，用于读取文件内容
	reader := bufio.NewReader(file)

	// 逐行读取文件内容
	for {
		line, err := reader.ReadString('\n') // 按照 '\n' 分割每一行内容
		if err != nil {
			// 判断文件是否已经读取完毕
			if err.Error() == "EOF" {
				break
			}
			utils.WriteLog("", err.Error())
		}

		var params connection.RedisParams
		err = json.Unmarshal([]byte(line), &params)
		if err != nil {
			fmt.Printf("Error: %v , content: %v", err, line)
		}

		wg.Add(1)
		goroutineChan <- struct{}{}
		go updateRedis(params)
	}

	//阻塞直到所有任务完成
	wg.Wait()
	for key, val := range dataInfo.DbToKey {
		//fmt.Printf("db number:%v\tkey count:%v\n", key, val)
		utils.WriteLog("", fmt.Sprintf("db number:%v\tkey count:%v", key, val))
		for k, v := range dataInfo.DbToKeyTypeToCount[key] {
			utils.WriteLog("", fmt.Sprintf("key type:%v\tkey count:%v", k, v))
		}
		fmt.Println()
	}
}

/**
 * updateRedis
 * @Description 更新redis
 * @Author Joe
 * @Date 2023-09-22 10:50:08
 * @Param data connection.RedisParams
 **/
func updateRedis(data connection.RedisParams) {
	conn := connection.RedisPool.Get()
	defer func() {
		<-goroutineChan
	}()
	defer wg.Done()
	defer conn.Close()

	conn.Do("select", data.Db)
	isSuccess := true

	switch data.Type {
	case "set":
		params := []interface{}{}
		params = append(params, data.Key)
		chunkParams := chunkArray(data.Value.([]interface{}), paramMaxCount)

		for _, v := range chunkParams {
			nextParam := append(params, v...)

			_, err := conn.Do("sadd", nextParam...)
			if err != nil {
				utils.WriteLog("", err.Error(), fmt.Sprintf("%v", params))
				isSuccess = false
			}
		}
	case "string":
		_, err := conn.Do("set", data.Key, data.Value)
		if err != nil {
			utils.WriteLog("", err.Error())
		}
	case "hash":
		params := []interface{}{}
		params = append(params, data.Key)

		values := []interface{}{}
		for k, v := range data.Value.(map[string]interface{}) {
			values = append(values, k, v)
		}

		chunkParams := chunkArray(values, paramMaxCount)

		for _, v := range chunkParams {
			nextParam := append(params, v...)

			_, err := conn.Do("hmset", nextParam...)
			if err != nil {
				utils.WriteLog("", err.Error(), fmt.Sprintf("%v", params))
				isSuccess = false
			}
		}
	case "zset":
		params := []interface{}{}
		params = append(params, data.Key)
		values := []interface{}{}
		for _, v := range data.Value.([]interface{}) {
			values = append(values, v.([]interface{})[1], v.([]interface{})[0])
		}
		chunkParams := chunkArray(values, paramMaxCount)

		for _, v := range chunkParams {
			nextParam := append(params, v...)

			_, err := conn.Do("zadd", nextParam...)
			if err != nil {
				utils.WriteLog("", err.Error(), fmt.Sprintf("%v", params))
				isSuccess = false
			}
		}
	case "list":
		params := []interface{}{}
		params = append(params, data.Key)
		chunkParams := chunkArray(data.Value.([]interface{}), paramMaxCount)

		for _, v := range chunkParams {
			nextParam := append(params, v...)

			_, err := conn.Do("rpush", nextParam...)
			if err != nil {
				utils.WriteLog("", err.Error(), fmt.Sprintf("%v", params))
				isSuccess = false
			}
		}
	}

	if isSuccess && data.TTL > 0 {
		conn.Do("EXPIRE", data.Key, data.TTL)
	}

	mutex.Lock()
	if _, ok := dataInfo.DbToKey[data.Db]; !ok {
		dataInfo.DbToKey[data.Db] = 1
		dataInfo.DbToKeyTypeToCount[data.Db] = map[string]int{}
	} else {
		dataInfo.DbToKey[data.Db] += 1
	}
	if _, ok := dataInfo.DbToKeyTypeToCount[data.Db][data.Type]; !ok {
		dataInfo.DbToKeyTypeToCount[data.Db][data.Type] = 1
	} else {
		dataInfo.DbToKeyTypeToCount[data.Db][data.Type] += 1
	}
	mutex.Unlock()
}

func chunkArray(arr []interface{}, count int) [][]interface{} {
	var result [][]interface{}
	size := len(arr)
	for i := 0; i < size; i += count {
		end := i + count
		if end > size {
			end = size
		}
		result = append(result, arr[i:end])
	}
	return result
}
