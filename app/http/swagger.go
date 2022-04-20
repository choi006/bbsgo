// Package http API.
// @title zhihu
// @version 1.0
// @description zhihu测试
// @termsOfService https://github.com/swaggo/swag

// @contact.name choi
// @contact.email choi

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

package http

import (
	_ "github.com/choi006/bbsgo/app/http/swagger"
)
