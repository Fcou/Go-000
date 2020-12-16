# ContentManagementSystem
# 内容管理系统（新闻发布管理系统）
#### 利用了Beego框架，基于MVC框架
![](https://mdn.mozillademos.org/files/16042/model-view-controller-light-blue.png)
#### M 对应本项目的models文件夹，数据库选用MySQL，操作库改用xorm，感觉比orm更加方便、简洁
- UserInfo 用户表：记录用户注册信息
- Article 文章表：记录每篇文章相关信息
- ArticleType 文章类型表：记录文章类型信息，与文章表是1：N的关系
- ArticleViewUser 访问记录表：记录访问信息，访问关系下，每篇文章与每个用户是多对多关系

#### V 对应本项目的views文件夹，包括各种页面html，相关静态资源例如css\img\js，存放在static文件夹下

#### C 对应本项目的controllers文件夹，包括各种页面处理逻辑
- article文件，主要包括 ArticleController 文章相关操作控制器，处理各种文章、文章类型相关操作
- user文件，主要包括 RegController 用户注册控制器，处理用户注册相关操作;LoginController 用户登录控制器，处理用户登录、登出相关操作

##### 更新：部分静态数据（文章类型）用Redis存储，提高效率

