package main

import (
	"github.com/hajimehoshi/ebiten"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer  FrameBuffer
	textures     TextureManager
	fizzleBuffer []int
}

func NewRenderer(width int, height int, manager TextureManager) Renderer {
	fb := NewFrameBuffer(width, height, black)
	fizzleBuffer := make([]int, 9999)
	for i := range fizzleBuffer {
		fizzleBuffer[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fizzleBuffer), func(i, j int) { fizzleBuffer[i], fizzleBuffer[j] = fizzleBuffer[j], fizzleBuffer[i] })
	//fmt.Printf("%v/n", fizzleBuffer)

	return Renderer{fb, manager, fizzleBuffer}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, screen *ebiten.Image) {
	width, height := screen.Size()

	focalLength := 1.0
	viewWidth := float64(width) / float64(height)
	viewPlane := world.player.direction.orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := MultiplyVector(world.player.direction, focalLength)
	viewCenter.Add(world.player.position)
	viewStart := DivideVector(viewPlane, 2)
	viewStart = SubstractVectors(viewCenter, viewStart)

	// sort sprites by distance from player, greatest distance first
	spritesByDistance := make(map[float64]Billboard)
	for _, sprite := range world.sprites(r.textures) {
		spriteDistance := SubstractVectors(sprite.start, world.player.position).length()
		spritesByDistance[spriteDistance] = sprite
	}
	keys := make([]float64, 0)
	for k, _ := range spritesByDistance {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	// Cast rays
	columns := width
	step := DivideVector(viewPlane, float64(columns))
	columnPosition := viewStart
	for x := 0; x < columns; x++ {
		rayDirection := SubstractVectors(columnPosition, world.player.position)
		viewPlaneDistance := rayDirection.length()
		ray := Ray{world.player.position, DivideVector(rayDirection, viewPlaneDistance)}
		lineEnd := world.worldmap.hitTest(ray)
		wallDistance := SubstractVectors(lineEnd, ray.origin).length()

		// Draw wall
		wallHeight := 1.0
		distanceRatio := viewPlaneDistance / focalLength
		perpendicular := wallDistance / distanceRatio
		realHeight := wallHeight * focalLength / perpendicular * float64(height)
		// is the wall a vertical (north/south) or horizontal (east/west)
		tile := world.worldmap.tile(lineEnd, ray.direction)
		wallX := lineEnd.x - math.Floor(lineEnd.x)
		wallTexture := r.textures.GetWallTextureByTile(tile, math.Floor(lineEnd.x) != lineEnd.x)
		if math.Floor(lineEnd.x) == lineEnd.x {
			wallX = lineEnd.y - math.Floor(lineEnd.y)
		}

		textureX := int(wallX * float64(wallTexture.image.Bounds().Size().X))
		// hack (substract a tiny ofset to prevent texture smearing)
		wallStart := Vector{float64(x), (float64(height)-realHeight)/2 - 0.001}
		r.frameBuffer.drawColumn(textureX, wallTexture, wallStart, realHeight, height)

		// Draw floor & ceiling
		floorStart := wallStart.y + float64(realHeight) + 1
		floorTile := Tile(-1)
		var floorTexture, ceilingTexture Texture
		for y := int(math.Min(floorStart, float64(height))); y < height; y++ {
			normalizedY := (float64(y)/float64(height))*2 - 1
			perpendicular := wallHeight * focalLength / normalizedY
			distance := perpendicular * distanceRatio
			mapPosition := MultiplyVector(ray.direction, distance)
			mapPosition.Add(ray.origin)
			tileX := math.Floor(mapPosition.x)
			tileY := math.Floor(mapPosition.y)
			tile := world.worldmap.GetTile(int(tileX), int(tileY))
			if tile != floorTile {
				floorTexture = r.textures.GetFloorCeilingTextureByTile(tile, false)
				ceilingTexture = r.textures.GetFloorCeilingTextureByTile(tile, true)
				floorTile = tile
			}

			textureX := mapPosition.x - tileX
			textureY := mapPosition.y - tileY
			r.frameBuffer.SetColorAt(x, y, floorTexture.GetColorAtNormalized(textureX, textureY))
			r.frameBuffer.SetColorAt(x, height-1-y, ceilingTexture.GetColorAtNormalized(textureX, textureY))
		}

		// Draw sprites
		for _, spriteDistance := range keys {

			if spriteDistance > wallDistance {
				continue
			}

			sprite := spritesByDistance[spriteDistance]

			hit := sprite.hitTest(ray)

			perpendicular = spriteDistance / distanceRatio
			spriteHeight := wallHeight / perpendicular * float64((height))
			spriteX := SubstractVectors(hit, sprite.start).length() / sprite.length
			spriteTexture := sprite.texture
			textureX = int(math.Min(spriteX*float64(spriteTexture.Width()), float64(spriteTexture.Width()-1)))
			start := Vector{float64(x), (float64(height)-spriteHeight)/2 + 0.001}
			r.frameBuffer.drawColumn(textureX, spriteTexture, start, spriteHeight, height)
		}

		columnPosition.Add(step)
	}

	// Effects
	for _, effect := range world.effects {
		switch effect.effectType {
		case fadeIn:
			r.frameBuffer.tint(effect.color, 1-effect.progress())
		case fadeOut:
			r.frameBuffer.tint(effect.color, effect.progress())
		case fizzleOut:
			threshold := int(effect.progress() * float64(len(r.fizzleBuffer)))
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					granularity := 2
					index := y/granularity*width + x/granularity
					fizzledIndex := r.fizzleBuffer[index%len(r.fizzleBuffer)]
					if fizzledIndex <= threshold {
						r.frameBuffer.SetColorAt(x, y, effect.color)
					}
				}
			}
		}
	}

	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}

func (r *Renderer) draw2d(world World, screen *ebiten.Image) {
	_, height := screen.Size()
	scale := float64(height) / float64(world.worldmap.Height)

	// Draw map
	for y := 0; y < world.worldmap.Height; y++ {
		for x := 0; x < world.worldmap.Width; x++ {
			if world.worldmap.GetTile(x, y).isWall() {
				rect := Rect{Vector{float64(x) * scale, float64(y) * scale},
					Vector{float64((x + 1)) * scale, float64((y + 1)) * scale}}
				r.frameBuffer.Fill(rect, white)
			}
		}
	}

	// Draw player
	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	r.frameBuffer.Fill(rect, blue)

	// Draw view plane
	focalLength := 1.0
	viewWidth := 1.0
	viewPlane := world.player.direction.orthogonal()
	viewPlane.Multiply(viewWidth)
	viewCenter := MultiplyVector(world.player.direction, focalLength)
	viewCenter.Add(world.player.position)
	viewStart := DivideVector(viewPlane, 2)
	viewStart = SubstractVectors(viewCenter, viewStart)
	viewEnd := AddVectors(viewStart, viewPlane)
	viewEnd.Multiply(scale)
	r.frameBuffer.DrawLine(MultiplyVector(viewStart, scale), viewEnd, red)

	// Cast rays
	columns := 10.0
	step := DivideVector(viewPlane, columns)
	columnPosition := viewStart
	for i := 0; i < int(columns); i++ {
		rayDirection := SubstractVectors(columnPosition, world.player.position)
		viewPlaneDistance := rayDirection.length()
		ray := Ray{world.player.position, DivideVector(rayDirection, viewPlaneDistance)}
		end := world.worldmap.hitTest(ray)
		for _, sprite := range world.sprites(r.textures) {
			hit := sprite.hitTest(ray)
			if (hit == Vector{}) { // does not work for vector 0, 0???
				continue
			}
			spriteDistance := SubstractVectors(hit, ray.origin).length()
			if spriteDistance > (SubstractVectors(end, ray.origin).length()) {
				continue
			}
			end = hit
		}
		start := MultiplyVector(ray.origin, scale)
		end.Multiply(scale)
		r.frameBuffer.DrawLine(start, end, green)
		columnPosition.Add(step)
	}

	for _, line := range world.sprites(r.textures) {
		r.frameBuffer.DrawLine(MultiplyVector(line.start, scale), MultiplyVector(line.end, scale), green)
	}

	screen.DrawImage(r.frameBuffer.ToImage(), nil)
}
