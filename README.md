# redis-dump

General redis-dump options:<br>
  $\qquad$-auth string<br>
        $\qquad$$\qquad$redis auth<br>
        $\qquad$$\qquad$redis 密码<br>
  $\qquad$-debug<br>
        $\qquad$$\qquad$is debug<br>
        $\qquad$$\qquad$debug开关<br>
  $\qquad$-file string<br>
        $\qquad$$\qquad$file path<br>
        $\qquad$$\qquad$文件路径 (default "redis_data.json")<br>
  $\qquad$-gcount int<br>
        $\qquad$$\qquad$goroutine count by redis key<br>
        $\qquad$$\qquad$并发执行key操作的最大协程数量 (default 100)<br>
  $\qquad$-help<br>
        $\qquad$$\qquad$Show help<br>
        $\qquad$$\qquad$显示帮助文档<br>
  $\qquad$-ip string<br>
        $\qquad$$\qquad$redis ip<br>
        $\qquad$$\qquad$redis ip地址 (default "127.0.0.1")<br>
  $\qquad$-op string<br>
        $\qquad$$\qquad$operation(dump or update)<br>
        $\qquad$$\qquad$操作类型(dump备份)(update更新)<br>
  $\qquad$-port int<br>
        $\qquad$$\qquad$redis port<br>
        $\qquad$$\qquad$redis 端口 (default 6379)<br>


