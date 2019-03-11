#### 参数设置
- -Z 单个文件的最大大小，默认为50,单位为MB
- -J 同时存在的文件个数，默认为5
- -BPF 字符串格式，bpf过滤表达式
- -NI 设置捕获的网卡ip，优先于NN参数
- -NN 设置捕获的网卡
- -ps 每个网络包的最大大小 默认为 1024(Byte)

#### 特殊设定
如果NN设置为any，则捕获全部网卡
如果ps设置为0，则实际值为Int32.Max

#### 样例

```bash
#windows:
.\dump-rotate.exe -Z 3 -NI 172.20.10.12
#linux:
./dump-rotate -Z 10 -NN any -ps 2048
```

#### parameter setting
- -Z use to specify single File size Max MB default 50MB
- -J file total sizes  default is five
- -BPF string bpf filter
- -NI specify capture ip,prior to NN parameter
- -NN specify network interface， any means all interface
- -ps single network packet size default is 1024(Byte)

#### special setting
if u set ps to zero, then the packetSize limit will be Int32.Max

#### sample

```bash
#windows:
.\dump-rotate.exe -Z 3 -NI 172.20.10.12
#linux:
./dump-rotate -Z 10 -NN any -ps 2048
```
