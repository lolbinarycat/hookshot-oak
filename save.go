package main

import (
	"bytes"
	"encoding/json"
	//"zfmt"
	"os"
	//"github.com/oakmound/oak/v2/fileutil"
)

//var SaveFileName = "save.json"

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

	err = json.Unmarshal(data,p)
	if err != nil {
		return err
	}

	return nil
}

func (m BasicPlayerModule) UnmarshalJSON(b []byte) error {
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
