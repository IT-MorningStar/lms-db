## 1. 基本介绍

### 对于学习rosedb、Mysql做个总结，用Go实现一个B+树版本的k-v数据库。

 - 实现目标
   - B+数
   - 4k分页（Page）
   - buffer-pool
   - 自适应hash索引 
   - write ahead log
   - lru页过期刷盘策略
   - sync page 数据异步、同步刷盘
   - check point机制
   - meta data 