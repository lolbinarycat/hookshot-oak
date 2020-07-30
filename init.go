package main

import (
	"github.com/lolbinarycat/hookshot-oak/physobj"
	"github.com/oakmound/oak/v2/event"
)

func init() {
	physobj.BlockBinding =
		func(blk *physobj.Block)  event.Bindable {
			return func(_ int, _ interface{}) int {
				blk.DoCollision(func() {
					if blk.Held == false {
						blk.DoGravity()
						if blk.ActiColls.GroundHit {
							blk.Body.Delta.SetPos(0, 0)
						}
					}
				})
				return 0
			}
		}
}
