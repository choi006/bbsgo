package provider

import (
    "github.com/choi006/bbsgo/app/provider/user"
    "github.com/gohade/hade/framework"
    "github.com/gohade/hade/framework/contract"
    "github.com/gohade/hade/framework/provider/app"
    "github.com/gohade/hade/framework/provider/cache"
    "github.com/gohade/hade/framework/provider/config"
    "github.com/gohade/hade/framework/provider/env"
    "github.com/gohade/hade/framework/provider/log"
    "github.com/gohade/hade/framework/provider/orm"
    "github.com/gohade/hade/framework/provider/redis"
    "testing"
)

const (
    BasePath = "/mnt/e/GoProjects/bbsgo"
)

func Test_Provider(t *testing.T) {
    // 初始化服务容器
    container := framework.NewHadeContainer()
    // 绑定App服务提供者
    container.Bind(&app.HadeAppProvider{BaseFolder: BasePath})
    // 后续初始化需要绑定的服务提供者...
    container.Bind(&env.HadeEnvProvider{})
    container.Bind(&config.HadeConfigProvider{})
    container.Bind(&log.HadeLogServiceProvider{})
    container.Bind(&orm.GormProvider{})
    container.Bind(&redis.RedisProvider{})
    container.Bind(&cache.HadeCacheProvider{})

    ormService := container.MustMake(contract.ORMKey).(contract.ORMService)
    db, err := ormService.GetDB()
    if err != nil {
        t.Fatal(err)
    }
    if err := db.AutoMigrate(&user.User{}); err != nil {
        t.Fatal(err)
    }
}
