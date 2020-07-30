package physobj

import (
	"bytes"
	"encoding/json"
)

type JSONVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

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
