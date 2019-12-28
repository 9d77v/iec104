# iec104
iec104主站golang实现

# quickstart
修改.env中的从站地址

执行docker-compose up -d

## develop
修改.env中的从站地址

./run d 进入docker容器

./run client 运行example中的104主站程序

修改example/client/worker/worker.go来处理通过104协议收到的数据

## 104规约解析
遥信起始地址1H<=>1
遥测起始地址4001H<=>16385
遥脉起始地址6001H<=>24577

## 实现功能

1. 每15分钟进行一次总召唤，第一次触发为激活后。

2. 每15分钟进行一次电度总召唤，第一次触发为总召唤结束后。

3. 信号量解析    
 
   3.1. M_SP_NA_1=1   单点遥信

   3.2. M_DP_NA_1=3   双点遥信

   3.3. M_ME_NA_1=9   带品质描述的遥测

   3.4 M_ME_NC_1=13   浮点数遥测

   3.4 M_IT_NA_1=15   电度总量遥脉

   3.5. M_SP_TB_1=30  带7个字节短时标的单点遥信

