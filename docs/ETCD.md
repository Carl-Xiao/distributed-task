### 开发环境
win10 + etcd单机部署+mongodb单机

#### etcd使用

- 查看所有的key

`etcdctl.exe get /cron/jobs --prefix`


- 租约监听
`etcdctl.exe get /cron/jobs --prefix`
> 能查看到所有项目进程ID

### 正式环境 阿里云服务器