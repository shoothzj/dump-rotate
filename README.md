#### 参数设置
- -Z use to specify single File size Max MB
- -J file total sizes
- -BPF 满足bpf表达式的过滤器
- -NI 指定抓包的ip地址,优先于NN参数
- -NN 指定抓包的网卡， any代表全抓

#### 样例

```bash
.\dump-rotate.exe -Z 3 -NI 172.20.10.12
```