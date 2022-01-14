// +build doc

package main

import (
	_ "github.com/lmxdawn/wallet/docs"
)

func init() {
	isSwag = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
