# govnstat
统计每个月VPS使用的流量，特别是AWS之类，防止流量溢出产生额外的费用

20240325: 修改代码，只统计发送流量也就是出站流量，因为根据nodeseek网友实验，AWS Lightsail的套餐内流量实际只统计出站流量 https://www.nodeseek.com/post-60328-1
## 使用 ##
1. VPS上先安装vnstat服务
2. 在运行go版本vnstat，当流量超出最大值，自动关机，建议采用service方式运行
3. 运行参数
```
Usage of ./vnstat:
  -loop int
      多少分钟检查检查一次流量使用情况
    	how many minutes to check loop (default 5) 
  -max float
      最大允许使用的流量，单位GB
    	max gb traffic todo (default 999)
  -p string
      vnstat服务的参数
    	vnstat args (default "-m")
  -ver int
      vnstat的版本，低于2.0必须使用此参数 -ver 1
    	the version of vnstat (default 2)
```
## Release 可执行文件已在centos和ubuntu server测试 ##
