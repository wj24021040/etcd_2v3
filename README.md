# etcd_2v3
etcd_2v3 is a simple command line utility to copy data from v2 to v3

# Installation
- Install golang(>=1.8)
- Build the binary :
     cd etcd_2v3
     go get ./...

# Usage
etcd_2v3 etcdV2(ip:port) etcdV3(ip:port) key(/)
ex:
etcd_2v3 192.168.2.1:2379 192.168.2.2:2379 /coreos.com