package initialize

import (
	"tiktok/pkg/constants"
	"tiktok/pkg/tracer"
)

func Jaeger() {
	tracer.InitJaeger(constants.ApiServiceName)
}
