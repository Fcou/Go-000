# 学习笔记
#### 1. [第二周课程笔记](https://naotu.baidu.com/file/0a6efe646327d237f3cac2b9c81e215f)
#### 2. 本周问题：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
* 思路：
    * 1.在dao层，通常会使用xorm\gorm等第三方库操作MySQL，我们应该把利用Wrap把error抛给上层处理。
    ```
    如果和其他库进行协作，考虑使用 errors.Wrap 或者 errors.Wrapf 保存堆栈信息。
    ```
    * 2.在server层，如果不产生新的业务逻辑错误，我们直接向上返回dao层的错误，不做任何处理。如果产生了新的逻辑错误，则利用errors.New返回错误信息
    ```
    在你的应用代码中，使用 errors.New 或者  errors.Errorf 返回错误。
    直接返回错误，而不是每个错误产生的地方到处打日志。
    ```
    * 3.在control层，这是调用的顶层，则处理传递上来的error，并且不再向上传递。使用 errors.Cause 获取 root error，再进行和 sentinel error 判定，或者断言错误实现了特定的行为，做出对应处理。
    ```
    在程序的顶部或者是工作的 goroutine 顶部(请求入口)，使用 %+v 把堆栈详情记录。
    使用 errors.Cause 获取 root error，再进行和 sentinel error 判定。
    ```