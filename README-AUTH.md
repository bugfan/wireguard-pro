## 

### wg set 部分修改
```
// all bak.conf
private_key=606527c72b83849a95d377aab6598de75fa651d896e6815adca187a633190e52
listen_port=51820
fwmark=0
replace_peers=true
public_key=b0edefa7b11dbfa0bd1eb62d7a7a14e13e49b7fff1e4ac05779af122d8a26344
replace_allowed_ips=true
allowed_ip=10.0.0.10/32
public_key=e2428b443df5217904be4ae9f474af949f481d0ee6fbe74b1622f8e9f1fcf07d
replace_allowed_ips=true
allowed_ip=10.0.0.100/32

// go.conf
private_key=606527c72b83849a95d377aab6598de75fa651d896e6815adca187a633190e52
listen_port=51820
fwmark=0
replace_peers=true

// wg set wg0 peer sO3vp7Edv6C9HrYtenoU4T5Jt//x5KwFd5rxItiiY0Q= allowed-ips 10.0.0.10/32
public_key=b0edefa7b11dbfa0bd1eb62d7a7a14e13e49b7fff1e4ac05779af122d8a26344
replace_allowed_ips=true
allowed_ip=10.0.0.10/32

```