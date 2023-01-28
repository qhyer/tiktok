package censor

import (
	"github.com/StellarisW/go-sensitive"
	"github.com/cloudwego/kitex/pkg/klog"
)

var TextCensor *sensitive.Manager

func Init() {
	filterManager := sensitive.NewFilter(
		sensitive.StoreOption{
			Type: sensitive.StoreMemory,
		},
		sensitive.FilterOption{
			Type: sensitive.FilterDfa,
		},
	)

	// 增加词汇
	err := filterManager.GetStore().AddWord("敏感词")
	if err != nil {
		klog.Errorf("load text censor failed %v", err)
		return
	}
	TextCensor = filterManager
}
