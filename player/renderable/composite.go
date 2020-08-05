package renderable

import (
	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/render"
)

type ComPlayerR struct {
	*render.CompositeR
}

func LoadCom() (ComPlayerR,error) {
	var comR *render.CompositeR
	var baseSpr, eyesSpr *render.Sprite
	var err error
	baseSpr, err = render.LoadSprite("", "player_new.png")
	if err != nil {
		goto EC
	}
	eyesSpr, err = render.LoadSprite("", "eyes.png")
	if err != nil {
		goto EC
	}

	comR = render.NewCompositeR(baseSpr, eyesSpr)
	return ComPlayerR{comR}, nil
EC:
	return ComPlayerR{}, err
}

func (c ComPlayerR) SetEyeDir(d direction.Dir) {
	c.AddOffset(1, floatgeom.Point2{d.HCoeff(),d.VCoeff()})
}
