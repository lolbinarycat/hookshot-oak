module github.com/lolbinarycat/hookshot-oak

go 1.14

replace github.com/lolbinarycat/utils => /home/binarycat/go/src/github.com/lolbinarycat/utils

require (
	github.com/disintegration/gift v1.2.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.0 // indirect
	github.com/lafriks/go-tiled v0.1.0
	github.com/lolbinarycat/utils v0.0.0-00010101000000-000000000000
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oakmound/oak/v2 v2.3.4-0.20200625001801-7a48ec29af75
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	github.com/yobert/alsa v0.0.0-20200618200352-d079056f5370 // indirect
	golang.org/x/image v0.0.0-20200430140353-33d19683fad8 // indirect
	golang.org/x/mobile v0.0.0-20200329125638-4c31acba0007
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200327173247-9dae0f8f5775 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/oakmound/oak/v2 => ./oak
