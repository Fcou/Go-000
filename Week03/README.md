# 学习笔记
#### 1. [第三周课程笔记](https://naotu.baidu.com/file/0a6efe646327d237f3cac2b9c81e215f)
#### 2. 本周问题：基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
* 思路：
    * 1.利用 errgroup 将多个 http server 并发执行
    * 2.调用 context 的 cancel 函数后，<-ctx.Done()不再阻塞，调用 server.Shutdown 函数，进而管控并发执行的 http server的生命周期
    * 3.利用 channel 监听 linux signal 信号，channel 接受到信号后，调用 context 的 cancel 函数，关闭全部服务
    * 4.如果 g.Wait() 出现非空error，则也调用 context 的 cancel 函数，关闭全部服务
    * 4.http server 全部关闭完成后通知主 goroutine