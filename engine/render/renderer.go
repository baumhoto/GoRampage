package render

import (
	"math"
	"math/rand"
	"sort"
	"time"

	_asset "github.com/baumhoto/GoRampage/engine/asset"
	_consts "github.com/baumhoto/GoRampage/engine/consts"
	_core "github.com/baumhoto/GoRampage/engine/core"
	_entity "github.com/baumhoto/GoRampage/engine/entity"
	_map "github.com/baumhoto/GoRampage/engine/map"
	"github.com/hajimehoshi/ebiten"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer    FrameBuffer
	fizzleBuffer   []int
}

func NewRenderer() Renderer {
	fb := NewFrameBuffer(_consts.SCREEN_WIDTH, _consts.SCREEN_HEIGHT, _consts.BLACK)
	fizzleBuffer := make([]int, 9999)
	for i := range fizzleBuffer {
		fizzleBuffer[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fizzleBuffer), func(i, j int) { fizzleBuffer[i], fizzleBuffer[j] = fizzleBuffer[j], fizzleBuffer[i] })
	//fmt.Printf("%v/n", fizzleBuffer)

	return Renderer{fb, fizzleBuffer}
}

// draw renders the world into the window
func (r *Renderer) Draw(world _entity.World, screen *ebiten.Image, tm *_asset.TextureManager) {
	width, height := screen.Size()

	focalLength := 1.0
	viewWidth := float64(width) / float64(height)
	viewPlane := world.Player.Direction.Orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := _core.MultiplyVector(world.Player.Direction, focalLength)
	viewCenter.Add(world.Player.Position)
	viewStart := _core.DivideVector(viewPlane, 2)
	viewStart = _core.SubstractVectors(viewCenter, viewStart)

	// sort sprites by distance from player, greatest distance first
	spritesByDistance := make(map[float64]_asset.Billboard)
	for _, sprite := range world.Sprites(tm) {
		spriteDistance := _core.SubstractVectors(sprite.Start, world.Player.Position).Length()
		spritesByDistance[spriteDistance] = sprite
	}
	keys := make([]float64, 0)
	for k, _ := range spritesByDistance {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	// Cast rays
	columns := width
	step := _core.DivideVector(viewPlane, float64(columns))
	columnPosition := viewStart
	for x := 0; x < columns; x++ {
		rayDirection := _core.SubstractVectors(columnPosition, world.Player.Position)
		viewPlaneDistance := rayDirection.Length()
		ray := _core.Ray{world.Player.Position, _core.DivideVector(rayDirection, viewPlaneDistance)}
		lineEnd := world.Worldmap.HitTest(ray)
		wallDistance := _core.SubstractVectors(lineEnd, ray.Origin).Length()

		// Draw wall
		wallHeight := 1.0
		distanceRatio := viewPlaneDistance / focalLength
		perpendicular := wallDistance / distanceRatio
		realHeight := wallHeight * focalLength / perpendicular * float64(height)
		// is the wall a vertical (north/south) or horizontal (east/west)
		tile := world.Worldmap.Tile(lineEnd, ray.Direction)
		wallX := lineEnd.X - math.Floor(lineEnd.X)
		wallTexture := tm.GetWallTextureByTile(tile, math.Floor(lineEnd.X) != lineEnd.X)
		if math.Floor(lineEnd.X) == lineEnd.X {
			wallX = lineEnd.Y - math.Floor(lineEnd.Y)
		}

		textureX := int(wallX * float64(wallTexture.Image.Bounds().Size().X))
		// hack (substract a tiny ofset to prevent texture smearing)
		wallStart := _core.Vector{X: float64(x), Y: (float64(height)-realHeight)/2 - 0.001}
		r.frameBuffer.drawColumn(textureX, wallTexture, wallStart, realHeight, height)

		// Draw floor & ceiling
		floorStart := wallStart.Y + float64(realHeight) + 1
		floorTile := _map.Tile(-1)
		var floorTexture, ceilingTexture _asset.Texture
		for y := int(math.Min(floorStart, float64(height))); y < height; y++ {
			normalizedY := (float64(y)/float64(height))*2 - 1
			perpendicular := wallHeight * focalLength / normalizedY
			distance := perpendicular * distanceRatio
			mapPosition := _core.MultiplyVector(ray.Direction, distance)
			mapPosition.Add(ray.Origin)
			tileX := math.Floor(mapPosition.X)
			tileY := math.Floor(mapPosition.Y)
			tile := world.Worldmap.GetTile(int(tileX), int(tileY))
			if tile != floorTile {
				floorTexture = tm.GetFloorCeilingTextureByTile(tile, false)
				ceilingTexture = tm.GetFloorCeilingTextureByTile(tile, true)
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
			spriteX := _core.SubstractVectors(hit, sprite.Start).Length() / sprite.Length
			spriteTexture := sprite.Texture
			textureX = int(math.Min(spriteX*float64(spriteTexture.Width()), float64(spriteTexture.Width()-1)))
			start := _core.Vector{float64(x), (float64(height)-spriteHeight)/2 + 0.001}
			r.frameBuffer.drawColumn(textureX, spriteTexture, start, spriteHeight, height)
		}

		columnPosition.Add(step)
	}

	// Player weapon
	r.frameBuffer.drawImage(
		tm.Animations[world.Player.Animation].Texture(world.Player.AnimationTime),
		_core.Vector{
			X: float64(width)/2.0 - float64(height)/2.0,
			Y: 0,
		}, _core.Vector{
			X: float64(height),
			Y: float64(height),
		})

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
				rect := _core.Rect{_core.Vector{float64(x) * scale, float64(y) * scale},
					_core.Vector{float64((x + 1)) * scale, float64((y + 1)) * scale}}
				r.frameBuffer.Fill(rect, _consts.WHITE)
			}
		}
	}

	// Draw player
	rect := world.Player.Rect()
	rect.Min.Multiply(scale)
	rect.Max.Multiply(scale)
	r.frameBuffer.Fill(rect, _consts.BLUE)

	// Draw view plane
	focalLength := 1.0
	viewWidth := 1.0
	viewPlane := world.Player.Direction.Orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := _core.MultiplyVector(world.Player.Direction, focalLength)
	viewCenter.Add(world.Player.Position)
	viewStart := _core.DivideVector(viewPlane, 2)
	viewStart = _core.SubstractVectors(viewCenter, viewStart)
	viewEnd := _core.AddVectors(viewStart, viewPlane)
	viewEnd.Multiply(scale)
	r.frameBuffer.DrawLine(_core.MultiplyVector(viewStart, scale), viewEnd, _consts.RED)

	// Cast rays
	columns := 10.0
	step := _core.DivideVector(viewPlane, columns)
	columnPosition := viewStart
	for i := 0; i < int(columns); i++ {
		rayDirection := _core.SubstractVectors(columnPosition, world.Player.Position)
		viewPlaneDistance := rayDirection.Length()
		ray := _core.Ray{Origin: world.Player.Position, Direction: _core.DivideVector(rayDirection, viewPlaneDistance)}
		end := world.Worldmap.HitTest(ray)
		for _, sprite := range world.Sprites(&_asset.TextureManager{}) {
			hit := sprite.HitTest(ray)
			if (hit == _core.Vector{}) { // does not work for vector 0, 0???
				continue
			}
			spriteDistance := _core.SubstractVectors(hit, ray.Origin).Length()
			if spriteDistance > (_core.SubstractVectors(end, ray.Origin).Length()) {
				continue
			}
			end = hit
		}
		start := _core.MultiplyVector(ray.Origin, scale)
		end.Multiply(scale)
		r.frameBuffer.DrawLine(start, end, _consts.GREEN)
		columnPosition.Add(step)
	}

	for _, line := range world.Sprites(&_asset.TextureManager{}) {
		r.frameBuffer.DrawLine(_core.MultiplyVector(line.Start, scale), _core.MultiplyVector(line.End, scale), _consts.GREEN)
	}

	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}

func (r *Renderer) ResetFrameBuffer() {
	r.frameBuffer.resetFrameBuffer()
}
