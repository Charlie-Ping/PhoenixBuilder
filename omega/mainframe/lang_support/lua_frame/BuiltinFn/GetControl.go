package BuiltinFn

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"

	lua "github.com/yuin/gopher-lua"
)

type BuiltGameControler struct {
	*BuiltinFn
}

func (b *BuiltGameControler) BuiltFunc(L *lua.LState) int {
	GameControl := L.NewTable()
	L.SetField(GameControl, "SendWsCmd", L.NewFunction(b.SendWsCmd))
	L.SetField(GameControl, "SendCmdAndInvokeOnResponse", L.NewFunction(b.SendCmdAndInvokeOnResponse))
	L.SetField(GameControl, "SetOnParamMsg", L.NewFunction(b.SetOnParamMsg))
	L.SetField(GameControl, "GetPos", L.NewFunction(b.GetPos))
	L.SetField(GameControl, "SayTo", L.NewFunction(ActionToDecorator(L, b.mainframe.GetGameControl().SayTo)))
	L.SetField(GameControl, "RawSayTo", L.NewFunction(ActionToDecorator(L, b.mainframe.GetGameControl().RawSayTo)))
	L.SetField(GameControl, "TitleTo", L.NewFunction(ActionToDecorator(L, b.mainframe.GetGameControl().TitleTo)))
	L.SetField(GameControl, "SubtitleTo", L.NewFunction(ActionToDecorator(L, b.mainframe.GetGameControl().SubTitleTo)))
	L.SetField(GameControl, "ActionBarTo", L.NewFunction(ActionToDecorator(L, b.mainframe.GetGameControl().ActionBarTo)))

	// 等待说话
	L.Push(GameControl)
	return 1
}

/*
	func (b *BuiltGameControler) BuiltGameContrler(L *lua.LState) int {
		GameControl := L.NewTable()
		L.SetField(GameControl, "SendWsCmd", L.NewFunction(b.SendWsCmd))
		L.SetField(GameControl, "SendCmdAndInvokeOnResponse", L.NewFunction(b.SendCmdAndInvokeOnResponse))
		//等待说话
		L.SetField(GameControl, "SetOnParamMsg", L.NewFunction(b.SetOnParamMsg))
		L.Push(GameControl)
		return 1
	}
*/
func (b *BuiltGameControler) SendWsCmd(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.ArgError(1, "参数应该只有一个")
	}
	args := L.CheckString(1)
	b.OmegaFrame.MainFrame.GetGameControl().SendCmd(args)

	return 1
}
func (b *BuiltGameControler) SendCmdAndInvokeOnResponse(L *lua.LState) int {
	if L.GetTop() == 1 {
		args := L.CheckString(1)
		ch := make(chan bool)
		b.OmegaFrame.MainFrame.GetGameControl().SendCmdAndInvokeOnResponse(args, func(output *packet.CommandOutput) {
			cmdBack := L.NewTable()
			if output.SuccessCount > 0 {
				L.SetField(cmdBack, "Success", lua.LBool(true))
			} else {
				L.SetField(cmdBack, "Success", lua.LBool(false))
			}
			L.SetField(cmdBack, "outputmsg", lua.LString(fmt.Sprintf("%v", output.OutputMessages)))
			L.Push(cmdBack)
			ch <- true
		})
		<-ch
	} else {
		// fmt.Println("参数应该仅有一个")
		L.ArgError(1, "参数应该只有一个")
	}
	return 1
}
func (b *BuiltGameControler) SetOnParamMsg(L *lua.LState) int {
	if L.GetTop() == 1 {
		name := L.CheckString(1)
		ch := make(chan bool)
		b.OmegaFrame.MainFrame.GetGameControl().SetOnParamMsg(name, func(chat *defines.GameChat) (catch bool) {
			msg := ""
			for _, v := range chat.Msg {
				msg += v + " "
			}
			L.Push(lua.LString(msg))
			ch <- true
			return false
		})
		<-ch
	} else {
		// fmt.Println("参数应该仅有一个")
		L.ArgError(1, "参数应该只有一个")
	}
	return 1
}

// This method originally belonged to PlayerKit,
// but for the sake of development efficiency,
// it was directly implemented in GameCtrl here temporary.
func (b *BuiltGameControler) GetPos(L *lua.LState) int {
	if L.GetTop() != 2 {
		L.ArgError(1, "takes exactly 2 arguments")
	}
	name := L.CheckString(1)
	selector := L.CheckString(2)

	PlayerPosChan := b.OmegaFrame.MainFrame.GetGameControl().GetPlayerKit(name).GetPos(selector)
	OriginPos := <-PlayerPosChan

	pos := L.NewTable()

	pos.Append(lua.LNumber(OriginPos.X()))
	pos.Append(lua.LNumber(OriginPos.Y()))
	pos.Append(lua.LNumber(OriginPos.Z()))
	L.Push(pos)
	return 1
}

func (b *BuiltGameControler) HasPermission(L *lua.LState) int {
	if L.GetTop() != 2 {
		L.ArgError(1, "takes exactly 2 arguments")
	}
	target := L.CheckString(1)
	key := L.CheckString(2)
	hp := b.OmegaFrame.MainFrame.GetGameControl().GetPlayerKit(target).HasPermission(key)
	L.Push(lua.LBool(hp))
	return 0
}

func ActionToDecorator(L *lua.LState, action func(string, string)) func(*lua.LState) int {
	return func(L *lua.LState) int {
		if L.GetTop() != 2 {
			L.ArgError(1, "takes exactly 2 arguments")
		}
		target := L.CheckString(1)
		line := L.CheckString(2)
		action(target, line)
		return 0
	}

}
