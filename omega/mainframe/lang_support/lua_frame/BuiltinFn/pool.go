package BuiltinFn

import lua "github.com/yuin/gopher-lua"

type BuiltinFunctionDic map[string]BuiltinFunctioner
type BuiltinFunctioner interface {
	BuiltFunc(L *lua.LState) int
}

// 获取内置函数 方便注入
func (b *BuiltinFn) GetSkynetBuiltinFunction() BuiltinFunctionDic {
	return map[string]BuiltinFunctioner{
		"GetListener":       &BuiltListener{b},
		"GetControl":        &BuiltGameControler{b},
		"loadComponent":     &LoadSide{b},
		"DataControler":     &BuiltDataControler{b},
		"GetBackEnder":      &BuiltBackEnder{BuiltinFn: b},
		"GetFilerControler": &BuiltFileControler{BuiltinFn: b},
	}
}
