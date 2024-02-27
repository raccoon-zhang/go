# 项目ing

> 可以正常连接代理gpt
> 
> 注册登陆功能
> 
> redis登陆校验功能，redis采用一主二从模式
> 
> ocr图片文字提取功能（可能gpt4会做的更好，但目前先用ocr代替，ocr环境需要自己配置),语言包： [Traineddata Files](https://github.com/tesseract-ocr/tesseract/wiki/Data-Files)

# 使用

## step1:添加你自己的apiKey到根目录下

> **.**
> 
> ├── README.md
> 
> ├── apiKey
> 
> ├── **dbPool**
> 
> │   ├── dbPool.go
> 
> │   ├── go.mod
> 
> │   └── go.sum
> 
> ├── **docs**
> 
> │   ├── createDatabase.md
> 
> │   └── createRedis.md
> 
> ├── **ginweb**
> 
> │   ├── ginweb.go
> 
> │   ├── go.mod
> 
> │   └── go.sum
> 
> ├── go.work
> 
> ├── **gptChat**
> 
> │   ├── go.mod
> 
> │   └── gptChat.go
> 
> ├── **main**
> 
> │   ├── go.mod
> 
> │   ├── go.sum
> 
> │   └── main.go
> 
> ├── **redisConfig**
> 
> │   ├── redis-sentinel1.conf
> 
> │   ├── redis-sentinel2.conf
> 
> │   ├── redis.conf
> 
> │   ├── redis6380.conf
> 
> │   ├── redis6381.conf
> 
> │   ├── start_redis.sh
> 
> │   └── stop_redis.sh
> 
> ├── **templates**
> 
> │   ├── 404.html
> 
> │   ├── chat.html
> 
> │   ├── **css**
> 
> │   │   ├── chat.css
> 
> │   │   ├── login.css
> 
> │   │   ├── registe.css
> 
> │   │   └── stable.css
> 
> │   ├── favicon.ico
> 
> │   ├── **js**
> 
> │   │   └── stable.js
> 
> │   ├── login.html
> 
> │   └── registe.html
> 
> └── **tools**
> 
>     ├── go.mod
> 
>     └── tool.go

## step2: 进入main目录下启动程序

# 参数配置

> 端口号：8080
