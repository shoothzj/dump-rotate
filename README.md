#### parameter setting
- -Z use to specify single File size Max MB default 50MB
- -J file total sizes  default is five
- -BPF string bpf filter
- -NI specify capture ip,prior to NN parameter
- -NN specify network interface， any means all interface
- -ps single network packet size default is 1024(Byte)

#### sample

```bash
#windows:
.\dump-rotate.exe -Z 3 -NI 172.20.10.12
#linux:
./dump-rotate -Z 10 -NN any -ps 2048
```
