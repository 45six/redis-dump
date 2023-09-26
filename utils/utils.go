// Package utils
// @Author Joe
// @Date 2023-09-23 21:04:32
// @Description:
package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var Debug bool

var DataFile string

/**
 * writeLog
 * @Description 写入log文件
 * @Author Joe
 * @Date 2023-09-22 15:19:50
 * @Param fileName string
 * @Param content string
 **/
func WriteLog(fileName string, params ...string) {
	fmt.Println(params[0])

	if !Debug {
		return
	}

	content := strings.Join(params, "\n")

	pc, file, line, ok := runtime.Caller(1)

	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		content = fmt.Sprintf("Time: %v\nFile: %s\nFunc: %s\nLine: %d\nContent: %v\n---------------------------------------------------------", time.Now().Format("2006-01-02 15:04:05"), file, funcName, line, content)
	}

	if fileName == "" {
		fileName = "log"
	}
	logFile, _ := os.OpenFile(fileName+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer logFile.Close()
	ret, err := logFile.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("log文件写入失败：%v --- %v \n", ret, err)
	}
}

/**
 * writeData
 * @Description 写入数据
 * @Author Joe
 * @Date 2023-09-22 15:20:02
 * @Param path string
 * @Param fileName string
 * @Param content string
 **/
func WriteData(path string, fileName string, content string) {
	os.MkdirAll(path, os.ModePerm)
	file, _ := os.OpenFile(path+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer file.Close()
	file.WriteString(content)
}

/*CheckExit
 * @Description 检查程序退出
 * @Author Joe
 * @Date 2023-09-23 21:40:05
 */
func CheckExit() {
	if runtime.GOOS == "windows" {
		fmt.Print("按回车键结束...")
		var strEnd string
		fmt.Scanf("%s", &strEnd)
	}
}
