# lntop

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/edouardparis/lntop/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/edouardparis/lntop)](https://goreportcard.com/report/github.com/edouardparis/lntop)
[![Godoc](https://godoc.org/github.com/edouardparis/lntop?status.svg)](https://godoc.org/github.com/edouardparis/lntop)
[![tippin.me](https://badgen.net/badge/%E2%9A%A1%EF%B8%8Ftippin.me/@edouardparis/F0918E)](https://tippin.me/@edouardparis)

`lntop` is an interactive text-mode channels viewer for Unix systems.

 ![lntop-v0.0.0](http://paris.iiens.net/lntop-v0.0.0.png?)
 *lntop-v0.0.0*

## Install

Require the [go programming language](https://golang.org/) (version >= 1.11)
```
git clone git@github.com:edouardparis/lntop.git
cd lntop && export GO111MODULE=on && go install ./...
```

## Config

First time `lntop` is used a config file `.lntop/config.toml` is created
in the user home directory.
```toml
[logger]
type = "production"
dest = "/root/.lntop/lntop.log"

[network]
name = "lnd"
type = "lnd"
address = "//127.0.0.1:10009"
cert = "/root/.lnd/tls.cert"
macaroon = "/root/.lnd/data/chain/bitcoin/mainnet/admin.macaroon"
macaroon_timeout = 60
max_msg_recv_size = 52428800
conn_timeout = 1000000
pool_capacity = 3
```
Change macaroon path according to your network.
