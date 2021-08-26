# MySQL 实例深度巡检 #

----------
##  Introductory ##

     您是否遇到过如下困扰：
       1）某一天业务反馈插入数据失败，经过排查自增id的数据类型溢出了，严重影响业务，而表数据量又很大，操作时间长，导致领导层震荡，小兵遭殃
       2）当前业务库下有多少个表字符集不是utf8或utf8mb4？
       3）当前业务库中有多少表引擎是非innodb？
       4）当前业务库是否存在冗余索引，存在多少？
       5）当前还有多少表没有自增主键，多少库还在用存储过程或触发器等？
       6）一些影响性能或安全的动态参数是否被修改，导致重启后或压力大时造成故障？
       7）当前业务库中存在哪些大表、表空间多大，表平均行度多大，是否需要通知业务进行大表的数据处理
       8）在出现重大节日保障、节庆促销、日常巡检（早高峰、下班后）时我们如何了解当前业务库是否能抗住压力不出现问题，让领导放心，让媳妇开心？
       等等
     针对如上问题，我开始针对MySQL实例中特定的检测项（影响安全与性能）进行分层次梳理，制定相应检查项的检测方法，为了减少重复性工作，减少手动巡检所浪费的时间开始使用go进行编写解放双手的工具，做一名合格的农民工。
    
    1）目前支持的检测项功能：
      1.1 数据库环境检查
          由于还在考虑，数据库环境是否需要检查，该检查哪些项，暂时还未实现该功能，不过接口以预留

      1.2 数据库配置检查
        1.2.01 双一是否设置  （sync_binlog、innodb_flush_log_at_trx_commit）
        1.2.02 只读权限是否关闭（tx_read_only、transaction_read_only、innodb_read_only、read_only、super_read_only）
        1.2.03 binlog格式是否为row（binlog_format=row）
        1.2.04 server端字符集是否为utf8（character_set_server=utf8）
        1.2.05 默认的密码认证插件是否为mysql_native_passwor（default_authentication_plugin）
        1.2.06 默认存储引擎及临时表的存储引擎是否为innodb（default_storage_engine、default_tmp_storage_engine、internal_tmp_disk_storage_engine）
        1.2.07 innodb脏页刷盘方式是否为O_DIRECT（innodb_flush_method）
        1.2.08 是否开启了死锁检测（innodb_deadlock_detect）
        1.2.09 查询缓存是否关闭（query_cache_type）
        1.2.10 与从库中继日志相关的参数（relay_log_purge=on,relay_log_recovery=on）
        1.2.11 检查当前事务隔离级别（transaction_isolation=READ-COMMITTED，tx_isolation=READ-COMMITTED）
        1.2.12 检查当前数据库的时区（system_time_zone=CTS,time_zone=system）
        1.2.13 检查是否开启了主键索引和唯一索引的重复行校验（unique_checks=on）

      1.3 数据库性能检查
        1.3.01 检测binlog落盘时使用磁盘的利用率（巡检项：binlogDiskUsageRate）
        1.3.02 历史最大连接数占最大连接数限制的百分比（巡检项：historyConnectionMaxUsageRate）
        1.3.03 创建临时磁盘表使用率（巡检项：tmpDiskTableUsageRate）
        1.3.04 创建临时磁盘文件使用率（巡检项：tmpDiskfileUsageRate）
        1.3.05 innodb buffer pool使用率（巡检项：innodbBufferPoolUsageRate）
        1.3.06 当前innodb buffer pool中脏页比例（巡检项：innodbBufferPoolUsageRate）
        1.3.07 当前innodb buffer pool命中率（巡检项：innodbBufferPoolHitRate）
        1.3.08 数据库文件句柄使用率（巡检项：openFileUsageRate）
        1.3.09 数据库表缓存率（巡检项：openTableCacheUsageRate）
        1.3.10 数据库表缓存溢出率（巡检项：openTableCacheOverflowsUsageRate）
        1.3.11 数据库全表扫描占比率（巡检项：selectScanUsageRate）
        1.3.12 数据库join语句全表扫描占比率（巡检项：selectfullJoinScanUsageRate）
        1.3.13 数据库表自增主键使用率（判断有符号及无符号int类型）（巡检项：tableAutoPrimaryKeyUsageRate）
        1.3.14 单表行数大于500w，且平均行长大于10KB（巡检项：tableRows）
        1.3.15 单表大于6G，并且碎片率大于30%（巡检项：diskFragmentationRate）
        1.3.16 单表行数大于1000W，且表空间大于30G（巡检项：bigTable）
        1.3.17 检查一个星期内未更新的表（巡检项：coldTable）

      1.4 数据库基线检查
        1.4.01 检查表字符集（输出非utf8或utf8mb4的表）（巡检项：tableCharset）
        1.4.02 检查引擎不是innodb的表（巡检项：tableEngine）
        1.4.03 检查表是否有外键关联（巡检项：tableForeign）
        1.4.04 检查表是否有自增主键（巡检项：tableNoPrimaryKey）
        1.4.05 检查主键自增列是否为bigint（巡检项：tableAutoIncrement）
        1.4.05 检查表中是否存在大字段blob、text、varchar(8099)、timestamp数据类型（巡检项：tableBigColumns）
        1.4.06 检查索引列是否允许为空（巡检项：indexColumnIsNull）
        1.4.07 检查索引列是否建立在enum、set、blob、text类型上（巡检项：indexColumnType）
        1.4.08 检查表中是否存在冗余索引（左匹配包含）（巡检项：tableIncludeRepeatIndex）
        1.4.09 检查库中是否存在存储过程、存储函数、触发器、视图（巡检项：tableProcedureFuncTriggerViews）

      1.5 数据库安全检查
        1.5.01 检查实例是否存在匿名用户（巡检项：anonymousUsers）
        1.5.02 检查实例是否存在空密码用户（巡检项：emptyPasswordUser）
        1.5.03 检查实例是否存在root用户是否允许远端登录（root@%）（巡检项：rootUserRemoteLogin）
        1.5.04 检查实例是否存在普通用户无访问ip限制（user@%）（巡检项：normalUserConnectionUnlimited）
        1.5.05 检查实例是否存在密码相同的用户（巡检项：userPasswordSame）
        1.5.06 检查实例是否存在普通用户所有库表权限（on *.*）（巡检项：normalUserDatabaseAllPrivilages）
        1.5.07 检查实例是否存在普通用户supr权限（WITH GRANT OPTION）（巡检项：normalUserSuperPrivilages）
        1.5.08 检查实例是否使用默认端口（3306）（巡检项：databasePort）

    2）mycheck下个版本改动(已完成)：
    （1）将配置参数检查外放，通过配置文件可进行自定义参数检测
    （2）将巡检项进行外放，通过配置文件可进行巡检项控制（通过false、true选择是否检测该项）
    （3）将巡检项阈值进行外放，通过配置文件可进行阈值调整（超过阈值才会出现在pdf巡检报告中）

    3）mycheck后续功能更新
    （1）增加数据库环境检查项
    （2）新添巡检结果输出方式，html格式
    （3）其他待续

    4）问题反馈
      如果在使用中有什么问题，或者有什么好的建议及想法，欢迎随时发送到邮箱ywlianghang@gmail.com或ywlianghang@163.com或xing.liang@greatdb.com
      注：如果新添巡检功能或巡检项，请写出巡检方法，巡检阈值及巡检意义

------

## Download  ##

   你可以从 [这里](https://github.com/ywlianghang/mycheck/releases) 下载二进制可执行文件，我已经在ubuntu、centos、redhat、windows x64下测试过

-----
## Usage  ##

   工具使用说明

    [root@lh-2 mycheck-linux-1.0]# ./mycheck -h
    NAME:
        mycheck - In-depth inspection of MySQL for system guarantee and daily inspection during major festivals

    USAGE:
        mycheck [global options] command [command options] [arguments...]

    VERSION:
        1.0

    AUTHOR:
        lianghang <ywlianghang@gmail.com>

    COMMANDS:
        help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
        --config value, -c value  Loading a Configuration File. for example： --config /tmp/mycheck.ymal (default: "nil")
        --help, -h                show help
        --version, -v             print the version

--------
## Examples ##

     1) 执行巡检命令方式一
     shell> ./mycheck --config mycheck.ymal
     2）执行巡检命令方式二
     shell> ./mycheck --config /tmp/mycheck.ymal
    
-------
## Building ##

    mycheck needs go version > 1.12 for go mod

    shell> git clone https://github.com/ywlianghang/mycheck.git
    shell> cd main
    shell> go build -o mycheck main.go
    shell> chmod +x mycheck
    shell> mv mycheck /usr/bin

-----
## Requirements ##

    巡检数据库实例必须为MySQL实例

-----
## Author ##

lianghang  ywlianghang@gmail.com
