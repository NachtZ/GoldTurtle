介绍
---
一个基于海龟交易法则，编写的纸黄金交易提示工具，在服务器测试运行半年，大概有20%的年化收益率。  
支持网页显示（9001端口），发送邮件提醒交易功能。  
目前存在bug，交易过程中存在某些计算问题，爬虫有时失败从而导致程序崩溃。
目前临时解决方案是用一个tutle.dat来存储关键参数，程序重启会自动读取该变量。
并利用supervisord托管程序，来保证高可用和临时的手动修正。

文件介绍
---

`turtle.pdf`:海龟交易法则pdf  
`db.go`:数据库操作  
`html.go`:网页显示  
`phaseXml.go`:海龟交易运算  
`runGold.go`:运行程序  
`addr.txt`:数据库地址  
`mail.txt`:邮箱地址  
`createTable.sql`:数据库建表  


未来计划
---
1. 清理这个版本的代码。
2. 修正目前存在的bug。  
3. 将程序改为td交易，增加收益和风险。  
4. 本程序还会收集每分钟的黄金数据，不知道该做什么用，考虑数据的利用。