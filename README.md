# go-search-pushbullet

for pushbullet重度使用者，结合chatbot使用，快速搜索出带有关键词的push

## 使用GOPROXY代理，加速拉取依赖的速度
如果你已经开始使用go 1.13，可以配置GOPROXY，加速拉取module依赖,关于module扶正，可以参见系列博文。
```shell
# goproxy.cn是国人提供的稳定代理，支持cdn加速
go env -w GOPROXY=https://goproxy.cn,direct
go env |grep -i proxy
```
实测下来，可以说是秒速拉取每个依赖