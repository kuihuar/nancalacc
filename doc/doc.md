snsapi_auth
  snsapi_login

snsapi_login
snsapi_auth

原来的
https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=dinglz1setxqhrpp7aa0&scope=snsapi_auth


https://oapi.dingtalk.com/connect/oauth2/sns_authorize?appid=SuiteKey&response_type=code&scope=snsapi_login&state=STATE&redirect_uri=REDIRECT_URI


https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=dinglz1setxqhrpp7aa0&scope=openid&prompt=consent




### 部署URL
[部署URL](https://365.kdocs.cn/l/cr165QH1rYos?openfrom=docs)


http://119.3.173.229/ksops/#/login

ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4	


ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4


mysql -h 10.27.10.225 -uwps -p'ffM48Cba86CfZAO5SBMV6pbPNN3HhTaKL-' -P3306


ffM48Cba86CfZAO5SBMV6pbPNN3HhTaKL-


    ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4	



    ecis-ecisaccountsync_db-168667383ff446ac80eb9fa93f2f3b33






    常用scope列表
基础权限

snsapi_login - 仅获取用户身份标识（unionid和openid），不需要用户确认

用户信息权限

snsapi_auth - 获取用户身份和授权信息

snsapi_userinfo - 获取用户身份和详细信息（需要用户确认）

企业相关权限

snsapi_org - 获取用户所在企业信息

snsapi_org_user - 获取用户在企业中的详细信息

使用建议
仅需要识别用户身份时，使用snsapi_login

需要获取用户详细信息时，使用snsapi_userinfo

需要获取企业相关信息时，添加企业相关scope

<!-- 示例授权URL
text
https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=YOUR_APPID&redirect_uri=YOUR_CALLBACK_URL&scope=snsapi_login,snsapi_userinfo&state=STATE


你现在单点登录配置完毕后， 默认系统后台进不去了， 可以通过 以下方式进入：
1）   http://119.3.173.229/account?oauthDumpDisable=true
这个访问这个链接， 然后通过 adminroot 以及对应的密码登录后， 在浏览器在换成 管理后台的地址 http://119.3.173.229/micsweb/sys/login  就可以再进入管理后台了

curl -X POST http://119.3.173.229/c/asyncacc/v1/account -H "Content-Type: application/json" -d '{"trigger_type":"1","sync_type":"1"}'
{"taskId":"20250716051038","createTime":"2025-07-16T05:10:37.591762493Z"} -->


https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=dinglz1setxqhrpp7aa0&scope=openid&prompt=consent
redirect_uri
state
1
1
1
1
1
code
http://119.3.173.229/c/asyncacc/v1/oauth/userAccessToken
http://119.3.173.229/c/asyncacc/v1/oauth/userinfo/me
userId
### 问题：

1. 多租户模式如何理解？
2. 中间表清理
3. 回滚同步（提单）
4. 增量同步接口触发
  - 增量父你部门问题（提单）
5. 原生API
5. 抛开插件

中间表，自己设计（）

用365文档实现



http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/incremen

nancal666


67da561345e9c9302a0ec2cd@default@ssh@root@172.16.170.22@22


67da561345e9c9302a0ec2cd@default@ssh@root@172.16.170.22@22@192.168.100.24


xI#24R838004

scp 


192.168.100.24


scp -p 2222 precheck_kubewps_dc_v7.0.2505b.20250610.254.tar.gz 67da561345e9c9302a0ec2cd@default@ssh@root@172.16.170.22@22@192.168.100.24

  <!-- - 删除部门，如果处理用户所胡在的部门 -->




各类平台登录账户信息如下：[kubewps内执行wpscli tool platform logininfo可再次查看]  
 所属平台    访问地址     管理员账号     管理员初始密码     
 ----------------- ----------------------------   ------------   -----------------  
 运维平台/部署平台 http://192.168.1.141/ksops/  wpsadmin     DGAW5ahouzw6$gxN1 
 运维平台/部署平台 https://192.168.1.141:8081/ksops/  wpsadmin     DGAW5ahouzw6$gxN 
 运维平台/部署平台 http://192.168.1.141/ksops/  guest  efK!1EYp9#jA!UCx 
 运维平台/部署平台 https://192.168.1.141:8081/ksops/  guest  efK!1EYp9#jA!UCx 
 v7文档中心管理后台 http://192.168.1.141   admin  EQ9s!8IiL@oU#KaO1
 v7文档中心系统后台 http://192.168.1.141/micsweb/sys/login   adminroot    un!mSnBOuv5IMaxs
 v7二开应用管理平台 http://192.168.1.141/c/manage/#/login    wpsadmin     QuLhnCuQepnE7o&3
 v7二开应用管理平台 http://192.168.1.141/c/manage/#/login    secadmin     ryX2VrYz65!mDDZn
 v7二开应用管理平台 http://192.168.1.141/c/manage/#/login    auditadmin   LdsKjm#j!3i8YC44


 https://deepwiki.com/



 navicat://conn.redis?Conn.AuthenticationMode=None&Conn.Host=localhost&Conn.Name=&Conn.Port=6379&Conn.SSH.AuthenticationMethod=Password&Conn.SSH.Host=192.168.1.141&Conn.SSH.Port=22&Conn.SSH.Username=root&Conn.Type=Standalone&Conn.UseSSH=true&Conn.UseSSL=false








 | 运维平台/部署平台  |http://119.3.173.229/ksops/ |  wpsadmin  | azcU@RJSiI19v!oz |   
| 运维平台/部署平台  |http://119.3.173.229/ksops/ |   guest    | 73zGSnA!Lkw8Meyq |   
| v7文档中心管理后台 |    http://119.3.173.229/account?oauthDumpDisable=true   |   admin    | awu&r3aN52!tsDkj |
| v7文档中心系统后台 | http://119.3.173.229/micsweb/sys/login | adminroot  | B5blVJ&2gESq37b! | 
| v7二开应用管理平台 | http://119.3.173.229/c/manage/#/login  |  wpsadmin  | a&6E&V#&1F2Vjwud 
| v7二开应用管理平台 | http://119.3.173.229/c/manage/#/login  |  secadmin  | m13U#AI0ePrq#@M0 |
| v7二开应用管理平台 | http://119.3.173.229/c/manage/#/login  | auditadmin | Lj9&R5rgK60LXS!@ 
| weboffice运维平台  | -  |     -|  -   | 可通过ksops运维平台顶部【业务导航】进行免密跳转访问！ |
|    woa管理后台     | -  |     -|  -   | 可通过ksops运维平台顶部【业务导航】进行免密跳转访问！ |



练习环境：
机器码： 8D0D4D56-EB5D-4465-496C-B611691ACAAF
私网IP： http://192.168.1.141/
服务器配置 16C48G
云文档编号：release_dc_v7.0.2505b.20250610-2025-06-04 18:45:51 CST-amd64



练习环境：
私网IP： http://192.168.1.141/
服务器配置 16C48G
云文档编号：release_dc_v7.0.2505b.20250610-2025-06-04 18:45:51 CST-amd64


39ml85pu4372.vicp.fun




Connect TiDB:    mysql --comments --host 127.0.0.1 --port 4000 -u root
TiDB Dashboard:  http://127.0.0.1:42547/dashboard
Grafana:         http://127.0.0.1:3000


curl --proto '=https' --tlsv1.2 -sSf https://tiup-mirrors.pingcap.com/install.sh | sh


tidb:
tiup playground v8.0.0  --db 1 --pd 1 --kv 1 --db.port 4000 --pd.port 4001 --host=192.168.1.142

tiup playground v8.0.0  --db 1 --pd 1 --kv 1 --db.port 4000 --pd.port 4001 --host 192.168.1.142
mongodb admin nopass