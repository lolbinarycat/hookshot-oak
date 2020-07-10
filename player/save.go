package player

import (
	"bytes"
	"encoding/json"
	//"zfmt"
	"os"
	//"github.com/oakmound/oak/v2/fileutil"
	"github.com/pkg/errors"
)

//var SaveFileName = "save.json"
//type JSONSave struct {
//	Player JSONPlayer
//}
type JSONPlayer struct {
	Pos,RespawnPos JSONVector
	Ctrls ControlConfig
	Mods PlayerModuleList
}
type JSONModList map[string]JSONMod
type JSONMod struct {
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
type JSONVector struct {
	X float64
	Y float64
}

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

/*func (m BasicPlayerModule) UnmarshalJSON(b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch tok {
		case nil:
			//noop
		case "Equipped":
			val, err := dec.Token()
			if err != nil {
				return err
			}
			if val == true {
				m.Obtain()
				m.Equip()
			}
		case "Obtained":
			val, err := dec.Token()
			if err != nil {
				return err
			}
			if val == true {
				m.Obtain()
			}
		}
	}
	return nil
}

func (m CtrldPlayerModule) UnmarshalJSON(b []byte) error {
	err := m.BasicPlayerModule.UnmarshalJSON(b)

	return err
}*/

func (m CtrldPlayerModule) UnmarshalJSON(b []byte) error {
	return PlayerModUnmarshalJSON(&m,b)
}

func (m CtrldPlayerModule) MarshalJSON() ([]byte,error)  {
	return PlayerModMarshalJSON(&m)
}

func (m BasicPlayerModule) UnmarshalJSON(b []byte) error {
	return PlayerModUnmarshalJSON(&m,b)
}

func (m BasicPlayerModule) MarshalJSON() ([]byte, error) {
	return PlayerModMarshalJSON(&m)
}


func (l PlayerModuleList) UnmarshalJSON(b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))

	for dec.More() {
		toc, err := dec.Token()
		if err != nil {
			return err
		}
		if _, ok := toc.(json.Delim);ok {
			continue
		}

		err = dec.Decode(l[toc.(string)])
		if err != nil {
			return err
		}
	}
	return nil
}

//func (ac ActiveCollisions) UnmarshalJSON(b []byte) error {
//	return nil
//}

func (o PhysObject) MarshalJSON() ([]byte,error) {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	pos := JSONVector{X:o.Body.X(),Y:o.Body.Y()}
	enc.Encode(pos)
	return buf.Bytes(), nil
}

func (o *PhysObject) UnmarshalJSON(b []byte) error {
	rdr := bytes.NewReader(b)

	dec := json.NewDecoder(rdr)
	vec := JSONVector{}
	dec.Decode(&vec)

	o.Body.SetPos(vec.X,vec.Y)
	return nil
}

func (p Player) MarshalJSON() ([]byte,error) {
	return json.Marshal(
		JSONPlayer{
			JSONVector{p.Body.X(),p.Body.Y()},
			JSONVector(p.RespawnPos),
			p.Ctrls,
			p.Mods,
		})
}

func (p *Player) UnmarshalJSON(b []byte) error {
	jsonP :=  newJSONPlayer()
	err := json.Unmarshal(b,&jsonP)
	if err != nil {
		return errors.Wrapf(err,"Player.UnmarshalJSON: json.Unmarshal([]byte,%v)",jsonP)
	}
	p.Body.SetPos(jsonP.Pos.X,jsonP.Pos.Y)
	p.RespawnPos = Pos(jsonP.RespawnPos)
	p.Ctrls = jsonP.Ctrls
	p.Mods = jsonP.Mods
	return nil
}

func PlayerModMarshalJSON(m PlayerModule) ([]byte,error) {
	jsonM := JSONMod{}
	if ctrldM, ok := m.(*CtrldPlayerModule); ok {
		jsonM.Type = CtrldMod
		jsonM.InputNum = ctrldM.GetInputNum()
	} else {
		jsonM.Type = BasicMod
		jsonM.InputNum = -1
	}
	basicM := m.GetBasic()
	jsonM.Equipped = basicM.Equipped
	jsonM.Obtained = basicM.Obtained

	b, err := json.Marshal(jsonM)
	if err != nil {
		return []byte{}, errors.Wrap(err,"PlayerModMarshalJSON")
	}
	return b, nil
}

func PlayerModUnmarshalJSON(m PlayerModule,b []byte) error {
	jsonM := JSONMod{}
	err := json.Unmarshal(b,&jsonM)
	if err != nil {
		return errors.Wrapf(err,"PlayerModUnmarshalJSON(%v,%v) failed",m,b)
	}
	if jsonM.Obtained {
		m.Obtain()
	}
	if jsonM.Equipped {
		m.Equip()
	}
	if jsonM.Type == CtrldMod {
		ctrldM := m.(*CtrldPlayerModule)
		ctrldM.Bind(nil,jsonM.InputNum)
	}
	return nil
}

func newJSONPlayer() JSONPlayer {
	jplr := JSONPlayer{}
	plr := Player{}

	InitMods(&plr)
	jplr.Mods = plr.Mods
	jplr.Ctrls = plr.Ctrls
	return jplr
}
