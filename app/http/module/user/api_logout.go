package user

import (
    "fmt"
    "github.com/choi006/bbsgo/app/http/middleware/auth"
    provider "github.com/choi006/bbsgo/app/provider/user"
    "github.com/gohade/hade/framework/contract"
    "github.com/gohade/hade/framework/gin"
    "time"
)

// Logout 代表登出
// @Summary 用户登出
// @Description 调用表示用户登出
// @Accept  json
// @Produce  json
// @Tags user
// @Success 200 {string} Message "用户登出成功"
// @Security ApiKeyAuth
// @Router /user/logout [get]
func (api *UserApi) Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.ISetStatus(500).IText("用户未登录")
        return
    }
    claim, err := auth.ParseToken(c, token)
    if err != nil {
        c.ISetStatus(500).IText(err.Error())
        return
    }
    user := &provider.User{
        ID:       claim.ID,
        UserName: claim.Username,
    }

    if claim.StandardClaims.ExpiresAt < time.Now().Unix() {
        c.ISetStatus(500).IText("该token已过期")
        return
    }
    // token剩余过期时间
    duration := claim.StandardClaims.ExpiresAt - time.Now().Unix()

    // 设置token黑名单
    cacheService := c.MustMake(contract.CacheKey).(contract.CacheService)
    key := fmt.Sprintf("user:logout:%v", token)
    if err := cacheService.SetObj(c, key, user, time.Duration(duration)*time.Second); err != nil {
        c.ISetStatus(500).IText(err.Error())
        return
    }

    c.ISetOkStatus().IText("用户登出成功")
    return
}
