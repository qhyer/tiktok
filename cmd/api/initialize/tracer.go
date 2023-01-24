package initialize

import (
	"tiktok/pkg/constants"
	"tiktok/pkg/tracer"
)

func InitJaeger() {
	tracer.InitJaeger(constants.ApiServiceName)
}
