## 在docker环境下运行（推荐）

```shell
docker run -d --name gotutu -p 3095:3095 \
--add-host="host.docker.internal:host-gateway" \
--restart on-failure:5 \
-v [你的本地路径]:/opt/data \
gmeroo/gotutu:latest
```

注意这里的 `--add-host` 只需在linux系统下添加，mac或者windows系统无需添加该指令

host.docker.internal 指向你的本地主机的地址，在容器内访问本地主机时请使用host.docker.internal来
访问，因为localhost指向的是docker容器内部

请将 `[你的本地路径]` 替换为你主机上的路径，这是用来存放图片和配置项的