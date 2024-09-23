# Kratos Project Template

## project init

### some env
```
export GOPROXY=https://goproxy.cn
#export GOPRIVATE=*.gitee.com
#export GOSUMDB="off"
#git config --global url."ssh://git@git.gitee.com/".insteadOf https://git.gitee.com/
```

### Install Kratos
```
go install github.com/go-kratos/kratos/v2
```

### Create a template project
```
kratos new git.gitee.com/group_name/{xxx}
cd {xxx}
```

### tool init
```
make init
```

### new a new config
```
make initConfig
```



## run
```
// for dev environment,you can run in Dockfile in production environment.
make runApi         // run API service
make runConsumer    // run message queue consumer
make runScript      // run a demo deamon once
```



## tool
### after modifying config file
```
cd {xxx}
make config
```


### after modifying wire file
```
go generate
```
