package render

import (
	_asset "github.com/baumhoto/go-rampage/engine/asset"
	_common "github.com/baumhoto/go-rampage/engine/common"
	_const "github.com/baumhoto/go-rampage/engine/consts"
	_entity "github.com/baumhoto/go-rampage/engine/entity"
	_map "github.com/baumhoto/go-rampage/engine/map"
	"github.com/hajimehoshi/ebiten"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer  FrameBuffer
	textures     _asset.TextureManager
	fizzleBuffer []int
}

func NewRenderer() Renderer {
	fb := NewFrameBuffer(_const.SCREEN_WIDTH, _const.SCREEN_HEIGHT, _const.BLACK)
	fizzleBuffer := make([]int, 9999)
	for i := range fizzleBuffer {
		fizzleBuffer[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fizzleBuffer), func(i, j int) { fizzleBuffer[i], fizzleBuffer[j] = fizzleBuffer[j], fizzleBuffer[i] })
	//fmt.Printf("%v/n", fizzleBuffer)

	return Renderer{fb, *_asset.NewTextureManager(), fizzleBuffer}
}

// draw renders the world into the window
func (r *Renderer) Draw(world _entity.World, screen *ebiten.Image) {
	width, height := screen.Size()

	focalLength := 1.0
	viewWidth := float64(width) / float64(height)
	viewPlane := world.Player.Direction.Orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := _common.MultiplyVector(world.Player.Direction, focalLength)
	viewCenter.Add(world.Player.Position)
	viewStart := _common.DivideVector(viewPlane, 2)
	viewStart = _common.SubstractVectors(viewCenter, viewStart)

	// sort sprites by distance from player, greatest distance first
	spritesByDistance := make(map[float64]_asset.Billboard)
	for _, sprite := range world.Sprites(r.textures) {
		spriteDistance := _common.SubstractVectors(sprite.Start, world.Player.Position).Length()
		spritesByDistance[spriteDistance] = sprite
	}
	keys := make([]float64, 0)
	for k, _ := range spritesByDistance {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	// Cast rays
	columns := width
	step := _common.DivideVector(viewPlane, float64(columns))
	columnPosition := viewStart
	for x := 0; x < columns; x++ {
		rayDirection := _common.SubstractVectors(columnPosition, world.Player.Position)
		viewPlaneDistance := rayDirection.Length()
		ray := _common.Ray{world.Player.Position, _common.DivideVector(rayDirection, viewPlaneDistance)}
		lineEnd := world.Worldmap.HitTest(ray)
		wallDistance := _common.SubstractVectors(lineEnd, ray.Origin).Length()

		// Draw wall
		wallHeight := 1.0
		distanceRatio := viewPlaneDistance / focalLength
		perpendicular := wallDistance / distanceRatio
		realHeight := wallHeight * focalLength / perpendicular * float64(height)
		// is the wall a vertical (north/south) or horizontal (east/west)
		tile := world.Worldmap.Tile(lineEnd, ray.Direction)
		wallX := lineEnd.X - math.Floor(lineEnd.X)
		wallTexture := r.textures.GetWallTextureByTile(tile, math.Floor(lineEnd.X) != lineEnd.X)
		if math.Floor(lineEnd.X) == lineEnd.X {
			wallX = lineEnd.Y - math.Floor(lineEnd.Y)
		}

		textureX := int(wallX * float64(wallTexture.Image.Bounds().Size().X))
		// hack (substract a tiny ofset to prevent texture smearing)
		wallStart := _common.Vector{float64(x), (float64(height)-realHeight)/2 - 0.001}
		r.frameBuffer.drawColumn(textureX, wallTexture, wallStart, realHeight, height)

		// Draw floor & ceiling
		floorStart := wallStart.Y + float64(realHeight) + 1
		floorTile := _map.Tile(-1)
		var floorTexture, ceilingTexture _asset.Texture
		for y := int(math.Min(floorStart, float64(height))); y < height; y++ {
			normalizedY := (float64(y)/float64(height))*2 - 1
			perpendicular := wallHeight * focalLength / normalizedY
			distance := perpendicular * distanceRatio
			mapPosition := _common.MultiplyVector(ray.Direction, distance)
			mapPosition.Add(ray.Origin)
			tileX := math.Floor(mapPosition.X)
			tileY := math.Floor(mapPosition.Y)
			tile := world.Worldmap.GetTile(int(tileX), int(tileY))
			if tile != floorTile {
				floorTexture = r.textures.GetFloorCeilingTextureByTile(tile, false)
				ceilingTexture = r.textures.GetFloorCeilingTextureByTile(tile, true)
				floorTile = tile
			}

			textureX := mapPosition.X - tileX
			textureY := mapPosition.Y - tileY
			r.frameBuffer.SetColorAt(x, y, floorTexture.GetColorAtNormalized(textureX, textureY))
			r.frameBuffer.SetColorAt(x, height-1-y, ceilingTexture.GetColorAtNormalized(textureX, textureY))
		}

		// Draw sprites
		for _, spriteDistance := range keys {

			if spriteDistance > wallDistance {
				continue
			}

			sprite := spritesByDistance[spriteDistance]

			hit := sprite.HitTest(ray)

			perpendicular = spriteDistance / distanceRatio
			spriteHeight := wallHeight / perpendicular * float64((height))
			spriteX := _common.SubstractVectors(hit, sprite.Start).Length() / sprite.Length
			spriteTexture := sprite.Texture
			textureX = int(math.Min(spriteX*float64(spriteTexture.Width()), float64(spriteTexture.Width()-1)))
			start := _common.Vector{float64(x), (float64(height)-spriteHeight)/2 + 0.001}
			r.frameBuffer.drawColumn(textureX, spriteTexture, start, spriteHeight, height)
		}

		columnPosition.Add(step)
	}

	// Effects
	for _, effect := range world.Effects {
		switch effect.EffectType {
		case _entity.FadeIn:
			r.frameBuffer.tint(effect.Color, 1-effect.Progress())
		case _entity.FadeOut:
			r.frameBuffer.tint(effect.Color, effect.Progress())
		case _entity.FizzleOut:
			threshold := int(effect.Progress() * float64(len(r.fizzleBuffer)))
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					granularity := 2
					index := y/granularity*width + x/granularity
					fizzledIndex := r.fizzleBuffer[index%len(r.fizzleBuffer)]
					if fizzledIndex <= threshold {
						r.frameBuffer.SetColorAt(x, y, effect.Color)
					}
				}
			}
		}
	}

	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}

func (r *Renderer) Draw2d(world _entity.World, screen *ebiten.Image) {
	_, height := screen.Size()
	scale := float64(height) / float64(world.Worldmap.Height)

	// Draw map
	for y := 0; y < world.Worldmap.Height; y++ {
		for x := 0; x < world.Worldmap.Width; x++ {
			if world.Worldmap.GetTile(x, y).IsWall() {
				rect := _common.Rect{_common.Vector{float64(x) * scale, float64(y) * scale},
					_common.Vector{float64((x + 1)) * scale, float64((y + 1)) * scale}}
				r.frameBuffer.Fill(rect, _const.WHITE)
			}
		}
	}

	// Draw player
	rect := world.Player.Rect()
	rect.Min.Multiply(scale)
	rect.Max.Multiply(scale)
	r.frameBuffer.Fill(rect, _const.BLUE)

	// Draw view plane
	focalLength := 1.0
	viewWidth := 1.0
	viewPlane := world.Player.Direction.Orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := _common.MultiplyVector(world.Player.Direction, focalLength)
	viewCenter.Add(world.Player.Position)
	viewStart := _common.DivideVector(viewPlane, 2)
	viewStart = _common.SubstractVectors(viewCenter, viewStart)
	viewEnd := _common.AddVectors(viewStart, viewPlane)
	viewEnd.Multiply(scale)
	r.frameBuffer.DrawLine(_common.MultiplyVector(viewStart, scale), viewEnd, _const.RED)

	// Cast rays
	columns := 10.0
	step := _common.DivideVector(viewPlane, columns)
	columnPosition := viewStart
	for i := 0; i < int(columns); i++ {
		rayDirection := _common.SubstractVectors(columnPosition, world.Player.Position)
		viewPlaneDistance := rayDirection.Length()
		ray := _common.Ray{world.Player.Position, _common.DivideVector(rayDirection, viewPlaneDistance)}
		end := world.Worldmap.HitTest(ray)
		for _, sprite := range world.Sprites(r.textures) {
			hit := sprite.HitTest(ray)
			if (hit == _common.Vector{}) { // does not work for vector 0, 0???
				continue
			}
			spriteDistance := _common.SubstractVectors(hit, ray.Origin).Length()
			if spriteDistance > (_common.SubstractVectors(end, ray.Origin).Length()) {
				continue
			}
			end = hit
		}
		start := _common.MultiplyVector(ray.Origin, scale)
		end.Multiply(scale)
		r.frameBuffer.DrawLine(start, end, _const.GREEN)
		columnPosition.Add(step)
	}

	for _, line := range world.Sprites(r.textures) {
		r.frameBuffer.DrawLine(_common.MultiplyVector(line.Start, scale), _common.MultiplyVector(line.End, scale), _const.GREEN)
	}

	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}

func (r *Renderer) ResetFrameBuffer() {
	r.frameBuffer.resetFrameBuffer()
}
