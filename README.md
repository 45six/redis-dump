# redis-dump

General redis-dump options:
  -auth string
        redis auth
        redis 密码
  -debug
        is debug
        debug开关
  -file string
        file path
        文件路径 (default "redis_data.json")
  -gcount int
        goroutine count by redis key
        并发执行key操作的最大协程数量 (default 100)
  -help
        Show help
        显示帮助文档
  -ip string
        redis ip
        redis ip地址 (default "127.0.0.1")
  -op string
        operation(dump or update)
        操作类型(dump备份)(update更新)
  -port int
        redis port
        redis 端口 (default 6379)
