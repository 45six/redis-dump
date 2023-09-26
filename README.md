# redis-dump ——redis 备份、更新工具

General redis-dump options:<br>
  &emsp;&emsp;-auth string<br>
        &emsp;&emsp;&emsp;&emsp;redis auth<br>
        &emsp;&emsp;&emsp;&emsp;redis 密码 (default "")<br>
  &emsp;&emsp;-debug<br>
        &emsp;&emsp;&emsp;&emsp;is debug<br>
        &emsp;&emsp;&emsp;&emsp;debug开关<br>
  &emsp;&emsp;-file string<br>
        &emsp;&emsp;&emsp;&emsp;file path<br>
        &emsp;&emsp;&emsp;&emsp;文件路径 (default "redis_data.json")<br>
  &emsp;&emsp;-gcount int<br>
        &emsp;&emsp;&emsp;&emsp;goroutine count by redis key<br>
        &emsp;&emsp;&emsp;&emsp;并发执行key操作的最大协程数量 (default 100)<br>
  &emsp;&emsp;-help<br>
        &emsp;&emsp;&emsp;&emsp;Show help<br>
        &emsp;&emsp;&emsp;&emsp;显示帮助文档<br>
  &emsp;&emsp;-ip string<br>
        &emsp;&emsp;&emsp;&emsp;redis ip<br>
        &emsp;&emsp;&emsp;&emsp;redis ip地址 (default "127.0.0.1")<br>
  &emsp;&emsp;-op string<br>
        &emsp;&emsp;&emsp;&emsp;operation(dump or update)<br>
        &emsp;&emsp;&emsp;&emsp;操作类型(dump备份)(update更新)<br>
  &emsp;&emsp;-port int<br>
        &emsp;&emsp;&emsp;&emsp;redis port<br>
        &emsp;&emsp;&emsp;&emsp;redis 端口 (default 6379)<br>


