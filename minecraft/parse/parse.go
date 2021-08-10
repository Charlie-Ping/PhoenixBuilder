package parse

import (
	"flag"
	"phoenixbuilder/minecraft/mctype"
	"strings"
	"fmt"
)

func Parse(Message string, defaultConfig *mctype.MainConfig) (*mctype.MainConfig, error) {
	//SLC := strings.Split(Message," ")
	isTransIMI:=false
	isInQuote:=false
	var SLC []string
	curmsg:=""
	for _,c:=range Message {
		if isTransIMI {
			isTransIMI=false
			curmsg+=string(c)
			continue
		}
		if(c=='\\') {
			isTransIMI=true
			continue
		}
		if(c=='"'){
			isInQuote=(!isInQuote)
			continue
		}
		if(c==' '&&!isInQuote){
			SLC=append(SLC,curmsg)
			curmsg=""
			continue
		}
		if(c=='#'){
			break
		}
		curmsg+=string(c)
	}
	if(len(curmsg)>0){
		SLC=append(SLC,curmsg)
		curmsg=""
	}
	//fmt.Printf("%v\n",SLC)
	if(isInQuote) {
		return nil, fmt.Errorf("Unterminated quoted string")
	}else if(isTransIMI) {
		return nil, fmt.Errorf("Unterminated escape")
	}
	Config := &mctype.MainConfig{
		Execute:   "",
		Block:     &mctype.ConstBlock{},
		OldBlock:  &mctype.ConstBlock{},
		//Begin:     mctype.Position{},
		End:       defaultConfig.End,
		Position:  defaultConfig.Position,
		Radius:    0,
		Length:    0,
		Width:     0,
		Height:    0,
		Method:    "replace",
		OldMethod: "keep",
	}

	FlagSet := flag.NewFlagSet("Parser", 0)
	var tempBlockData int
	var tempOldBlockData int
	//Length,  Width and Height
	FlagSet.IntVar(&Config.Length,"length",defaultConfig.Length,"The length")
	FlagSet.IntVar(&Config.Length,"l",defaultConfig.Length,"The length")
	FlagSet.IntVar(&Config.Width,"width",defaultConfig.Width,"The width")
	FlagSet.IntVar(&Config.Width,"w",defaultConfig.Width,"The width")
	FlagSet.IntVar(&Config.Height,"height",defaultConfig.Height,"The height")
	FlagSet.IntVar(&Config.Height,"h",defaultConfig.Height,"The height")
	//Radius
	FlagSet.IntVar(&Config.Radius,"radius",defaultConfig.Radius,"The radius")
	FlagSet.IntVar(&Config.Radius,"r",defaultConfig.Radius,"The radius")
	//Facing, Path, Shape
	FlagSet.StringVar(&Config.Facing,"facing",defaultConfig.Facing,"Building's facing")
	FlagSet.StringVar(&Config.Facing,"f",defaultConfig.Facing,"Building's facing")
	FlagSet.StringVar(&Config.Path,"path",defaultConfig.Path,"The path of file")
	FlagSet.StringVar(&Config.Path,"p",defaultConfig.Path,"The path of file")
	FlagSet.StringVar(&Config.Shape,"shape",defaultConfig.Shape,"The path of file")
	FlagSet.StringVar(&Config.Shape,"s",defaultConfig.Shape,"The path of file")
	//Block
	FlagSet.StringVar(&Config.Block.Name,"block",defaultConfig.Block.Name,"Blocks that make up the building")
	FlagSet.StringVar(&Config.Block.Name,"b",defaultConfig.Block.Name,"Blocks that make up the building")
	FlagSet.IntVar(&tempBlockData,"data",int(defaultConfig.Block.Data),"The data of Block")
	FlagSet.IntVar(&tempBlockData,"d",int(defaultConfig.Block.Data),"The data of Block")
	//OldBlock
	FlagSet.StringVar(&Config.OldBlock.Name,"old_block",defaultConfig.OldBlock.Name,"Blocks that make up the building")
	FlagSet.StringVar(&Config.OldBlock.Name,"ob",defaultConfig.OldBlock.Name,"Blocks that make up the building")
	FlagSet.IntVar(&tempOldBlockData,"old_data",int(defaultConfig.OldBlock.Data),"The data of Block")
	FlagSet.IntVar(&tempOldBlockData,"od",int(defaultConfig.OldBlock.Data),"The data of Block")

	FlagSet.Parse(SLC[1:])
	/*for k, _ := range builder.Builder {
		if k == SLC[0] {
			Config.Execute = k
		}
	}*/
	Config.Execute = SLC[0]
	// Since the function system exists ^^
	
	//for index, v := range SLC {
	//	if v == "-p" || v == "--position" {
	//		x, xe := strconv.Atoi(SLC[index + 1])
	//		y, ye := strconv.Atoi(SLC[index + 2])
	//		z, ze := strconv.Atoi(SLC[index + 3])
	//		if xe == nil && ye == nil && ze == nil {
	//			Config.Position = mctype.Position{X: x, Y: y, Z: z}
	//		}
	//	}
	//}
	Config.Block.Data=int16(tempBlockData)
	Config.OldBlock.Data=int16(tempOldBlockData)
	return Config, nil
}

func PipeParse(Message string, config *mctype.MainConfig) ([]*mctype.MainConfig, error) {
	ChatSlice := strings.Split(Message,"|")
	var Configs []*mctype.MainConfig
	for _, v := range ChatSlice {
		pv, err:=Parse(v,config)
		if err!=nil {
			return nil, err
		}
		Configs = append(Configs,pv)
	}
	return Configs, nil
}