# monster-go
2.etcd isWatch参数优化，和isGate有重合
3.output是否需要优化: 暂时没有太多优化
1.etcd,先启动gate，再启动login，login重启，etcd watch会先进put再进del，导致连接不上: 暂时在连接时做一定延时

buglist:

4.很多地方的日志问题，错误处理
center注册的是127.0.0.1的ip