package player

import (
	//"fmt"

	"encoding/json"
	"fmt"

	//"zfmt"
	"os"
	//"github.com/oakmound/oak/v2/fileutil"
	"github.com/pkg/errors"
	"github.com/lolbinarycat/hookshot-oak/physobj"
)

//type JSONSave struct {
//	Player JSONPlayer
//}
type JSONVector = physobj.JSONVector

type JSONPlayer struct {
	Pos,RespawnPos JSONVector
	Ctrls ControlConfig
	Mods []JSONMod
}
//type JSONModList map[string]JSONMod
type JSONMod struct {
	Name string
	// JSON doesn't encode types, so we do it manually with an enum.
	Type JSONModType
	Obtained bool
	Equipped bool
	InputNum int //0-7, or -1 if N/A
}
type JSONModType int
const (
	BasicMod JSONModType = iota
	CtrldMod
)


func (p Player) Save(saveFileName string) error {
	saveFile, err := os.Create(saveFileName)
	defer saveFile.Close()
	if err != nil {
		return err
	}
	saveData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = saveFile.Write(saveData)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) Load(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	fileStats, err := file.Stat()
	if err != nil {
		return err
	}
	data := make([]byte,fileStats.Size())
	file.Read(data)

	if p == nil {
		panic("p == nil")
	}
	err = json.Unmarshal(data,p)
	if err != nil {
		return errors.Wrapf(err,"%v.Load(%v) failed",p,filename)
	}

	return nil
}







//func (ac ActiveCollisions) UnmarshalJSON(b []byte) error {
//	return nil
//}


func (p Player) MarshalJSON() ([]byte,error) {
	var jMods = make([]JSONMod,len(p.Mods))
	var i int = 0
	for name, mod := range p.Mods {
		bMod := mod.GetBasic()
		modType := BasicMod
		inputNum := -1
		if ctrldM, isCtrld := mod.(*CtrldPlayerModule); isCtrld {
			modType = CtrldMod
			inputNum = ctrldM.GetInputNum()
		}
		jMods[i] = JSONMod{
			Name:name,
			Equipped:bMod.Equipped,
			Obtained:bMod.Obtained,
			Type: modType,
			InputNum: inputNum,
		}
		i++
	}
	return json.Marshal(
		JSONPlayer{
			JSONVector{p.Body.X(),p.Body.Y()},
			JSONVector(p.RespawnPos),
			p.Ctrls,
			jMods,
		})
}

func (p *Player) UnmarshalJSON(b []byte) error {
	jsonP := JSONPlayer{}
	jsonP.Ctrls = p.Ctrls
	err := json.Unmarshal(b,&jsonP)
	fmt.Println(jsonP.Ctrls.Mod)
	if err != nil {
		return errors.Wrapf(err,"Player.UnmarshalJSON: json.Unmarshal([]byte,%v)",jsonP)
	}
	p.Body.SetPos(jsonP.Pos.X,jsonP.Pos.Y)
	p.RespawnPos = Pos(jsonP.RespawnPos)
	p.Ctrls = jsonP.Ctrls
	for _, jMod := range jsonP.Mods {
		mn := jMod.Name // mod name
		switch jMod.Type {
		case BasicMod:
			p.Mods[mn] = &BasicPlayerModule{
				Equipped:jMod.Equipped,Obtained:jMod.Obtained}
		case CtrldMod:
			if jMod.Obtained {
				p.Mods[mn].Obtain()
				if jMod.Equipped {
					p.Mods[mn].Equip()
					p.Mods[mn].(*CtrldPlayerModule).Bind(p,jMod.InputNum)
				}
			}
		default:
			return errors.New("unknown json mod type")
		}
	}
	fmt.Println(p.Mods["jump"])
	return nil
}

