package builder

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"phoenixbuilder/minecraft/mctype"
)

type ColorBlock struct {
	Color colorful.Color
	Block *mctype.ConstBlock
}

func Paint(config *mctype.MainConfig, blc chan *mctype.Module) error {
	path := config.Path
	width := config.Width
	height := config.Height
	facing := config.Facing
	pos := config.Position
	img, err := imaging.Open(path)
	if err != nil {
		return err
	}
	if width != 0 && height != 0 {
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}
	Max := img.Bounds().Max
	X, Y := Max.X, Max.Y
	//BlockSet := make([]*mctype.Module, X*Y)
	index := 0
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := colorful.Color{
				R: float64(r & 0xff),
				G: float64(g & 0xff),
				B: float64(b & 0xff),
			}
			switch facing {
			default:
				return errors.New("Facing (-f) not defined")
			case "x":
				blc <- &mctype.Module{
					Point: mctype.Position{
						X: pos.X,
						Y: x + pos.Y,
						Z: y + pos.Z,
					},
					Block: getBlock(c),
				}
			case "y":
				blc <- &mctype.Module{
					Point: mctype.Position{
						X: x + pos.X,
						Y: pos.Y,
						Z: y + pos.Z,

					},
					Block: getBlock(c),
				}
			case "z":
				blc <- &mctype.Module{
					Point: mctype.Position{
						X: x + pos.X,
						Y: y + pos.Y,
						Z: pos.Z,
					},
					Block: getBlock(c),
				}
			}

			index++
		}
	}
	return nil
}

func getBlock(c colorful.Color) *mctype.Block {
	if _, _, _, a := c.RGBA(); a == 0 {
		return AirBlock.Take()
	}
	var List []float64
	for _, v := range ColorTable {
		s := c.DistanceRgb(v.Color)
		List = append(List, s)
	}
	return ColorTable[getMin(List)].Block.Take()
}

func getMin(t []float64) int {
	min := t[0]
	index := 0
	for i, v := range t {
		if v < min {
			min = v
			index = i
		}
	}
	return index
}

