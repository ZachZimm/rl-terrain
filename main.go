package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	tileSizeX = 1
	tileSizeY = 1
)

var (
	windowWidth  = 1200
	windowHeight = 675
	numTilesX    = windowWidth / tileSizeX
	numTilesY    = windowHeight / tileSizeY

	tileMap [1200][675]float64
	lehmer  *Lehmer
)

func initTileMap() {
	for i := 0; i < numTilesX; i++ {
		for j := 0; j < numTilesY; j++ {
			tileMap[i][j] = 0.5
		}
	}
}

func normalizeTileMap() {
	minValue := 1.0
	maxValue := 0.0
	for i := 0; i < numTilesX; i++ {
		for j := 0; j < numTilesY; j++ {
			minValue = math.Min(minValue, tileMap[i][j])
			maxValue = math.Max(maxValue, tileMap[i][j])
		}
	}

	for i := 0; i < numTilesX; i++ {
		for j := 0; j < numTilesY; j++ {
			tileMap[i][j] = (tileMap[i][j] - minValue) / (maxValue - minValue)
		}
	}
}

func setTilesRandomly_rl(willModify bool) {
	for i := 0; i < numTilesX; i++ {
		for j := 0; j < numTilesY; j++ {
			randValue := float64(rl.GetRandomValue(0, 100)) / 100.0
			if willModify {
				oldValue := tileMap[i][j]
				newValue := oldValue * randValue
				tileMap[i][j] = (oldValue + newValue) / 2
			} else {
				tileMap[i][j] = randValue
			}
		}
	}
}

func setTilesRandomly_rand(willModify bool) {
	for i := 0; i < numTilesX; i++ {
		for j := 0; j < numTilesY; j++ {
			randValue := rand.Float64()
			if willModify {
				oldValue := tileMap[i][j]
				newValue := oldValue * randValue
				tileMap[i][j] = (oldValue + newValue) / 2
			} else {
				tileMap[i][j] = randValue
			}
		}
	}
}

var (
	alpha       = 3.
	beta        = 4.
	n     int32 = 9
	seed  int64 = 100
)

func setTilesRandomly_perlin(willModify bool) {
	seed = rand.Int63()
	p := perlin.NewPerlin(alpha, beta, n, seed)
	for x := 0; x < numTilesX; x++ {
		for y := 0; y < numTilesY; y++ {
			randValue := p.Noise2D(float64(x)/float64(numTilesX), float64(y)/float64(numTilesY))
			randValue = math.Min(math.Max(randValue, 0), 1)
			if willModify {
				oldValue := tileMap[x][y]
				newValue := (oldValue * randValue) + 0.25
				newValue = math.Min(math.Max(newValue, 0), 1)
				tileMap[x][y] = (oldValue + newValue) / 2
			} else {
				tileMap[x][y] = randValue
			}
		}
	}
}

func setTilesRandomly_lehmer(willModify bool) {
	seed = lehmer.Int63()
	for x := 0; x < numTilesX; x++ {
		for y := 0; y < numTilesY; y++ {
			randInt := lehmer.Int63() % 100
			randValue := float64(randInt) / 100.0
			randValue = math.Min(math.Max(randValue, 0), 1)
			if willModify {
				oldValue := tileMap[x][y]
				newValue := (oldValue * randValue) + 0.25
				newValue = math.Min(math.Max(newValue, 0), 1)
				tileMap[x][y] = (oldValue + newValue) / 2
			} else {
				tileMap[x][y] = randValue
			}
		}
	}
}

func mapGenerationRoutine(willReset bool) {
	if willReset {
		initTileMap()
	}
	perlinSteps := 4
	for i := 0; i < perlinSteps; i++ {
		setTilesRandomly_perlin(true)
	}
	normalizeTileMap()
}

func getColorFromTile(tile float64) rl.Color {
	return rl.NewColor(0, uint8(tile*255), 0, 255)
}

func getColorFromTile_hsl(tile float64) rl.Color {
	hue := tile * 360
	return rl.ColorFromHSV(float32(hue), 0.5, 0.75)
}

// Deep water: 0.0 - 0.2
// Shallow water: 0.2 - 0.4
// Sand: 0.4 - 0.5
// Grass: 0.5 - 0.65
// Forest: 0.65 - 0.8
// Mountains: 0.8 - 0.9
// High mountains: 0.9 - 1.0

var deepWater = rl.NewColor(0, 0, 128, 255)
var shallowWater = rl.NewColor(0, 0, 255, 255)
var sand = rl.NewColor(240, 240, 64, 255)
var grass = rl.NewColor(0, 255, 0, 255)
var forest = rl.NewColor(0, 128, 0, 255)
var dirt = rl.NewColor(128, 64, 0, 255)
var mountains = rl.NewColor(128, 128, 128, 255)
var highMountains = rl.NewColor(255, 255, 255, 255)

func getColorFromSwitch(tile float64) rl.Color {
	// Find the bucket that the tile value falls into
	// using a switch statement
	switch {
	case tile < 0.18:
		return deepWater
	case tile < 0.3:
		return shallowWater
	case tile < 0.4:
		return sand
	case tile < 0.5:
		return grass
	case tile < 0.7:
		return forest
	case tile < 0.82:
		return dirt
	case tile < 0.95:
		return mountains
	default:
		return highMountains
	}
}

func main() {
	lehmer = NewLehmer(rand.Int63())
	initTileMap()
	const keyPressTimer = 2
	lastKeyPress := "1"
	lastKeyPressTime := time.Now()
	showLastKeyPress := true

	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Tilemap")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		// rl.ClearBackground(rl.RayWhite)

		for i := 0; i < numTilesX; i++ {
			for j := 0; j < numTilesY; j++ {
				rl.DrawRectangle(int32(i*tileSizeX), int32(j*tileSizeY), int32(tileSizeX), int32(tileSizeY), getColorFromSwitch(tileMap[i][j]))
			}
		}
		rl.DrawRectangle(5, int32(float32(windowHeight)*0.91), 320, 30, rl.Fade(rl.White, 0.5))
		if showLastKeyPress && time.Since(lastKeyPressTime).Seconds() < keyPressTimer {
			rl.DrawText("Last key pressed: "+lastKeyPress, 10, int32(float32(windowHeight)*0.92), 20, rl.Black)
		}

		rl.DrawRectangle(5, 5, 500, 170, rl.Fade(rl.White, 0.5))
		// Display the controls on the screen
		rl.DrawText("Press R to randomize (raylib)", 10, 10, 20, rl.Black)
		rl.DrawText("Press M to randomize (rand)", 10, 30, 20, rl.Black)
		rl.DrawText("Press P to randomize (perlin)", 10, 50, 20, rl.Black)
		rl.DrawText("Press L to randomize (lehmer)", 10, 70, 20, rl.Black)
		rl.DrawText("Press N to normalize", 10, 90, 20, rl.Black)
		rl.DrawText("Press G to generate with a layered approach", 10, 110, 20, rl.Black)
		rl.DrawText("Hold shift to layer (multiply) the random values", 10, 130, 20, rl.Black)
		rl.DrawText("Press 1 to reset", 10, 150, 20, rl.Black)

		rl.EndDrawing()
		if rl.IsKeyPressed(rl.KeyR) {
			lastKeyPress = "Shift + R"
			willModify := true
			if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
				lastKeyPress = "R"
				initTileMap()
				willModify = false
			}
			setTilesRandomly_rl(willModify)
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyM) {
			lastKeyPress = "Shift + M"
			willModify := true
			if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
				lastKeyPress = "M"
				initTileMap()
				willModify = false
			}
			setTilesRandomly_rand(willModify)
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyP) {
			lastKeyPress = "Shift + P"
			willModify := true
			if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
				lastKeyPress = "P"
				initTileMap()
				willModify = false
			}
			setTilesRandomly_perlin(willModify)
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyL) {
			lastKeyPress = "Shift + L"
			willModify := true
			if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
				lastKeyPress = "L"
				initTileMap()
				willModify = false
			}
			setTilesRandomly_lehmer(willModify)
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyN) {
			lastKeyPress = "N"
			normalizeTileMap()
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyOne) {
			lastKeyPress = "1"
			initTileMap()
			lastKeyPressTime = time.Now()
		} else if rl.IsKeyPressed(rl.KeyG) {
			lastKeyPress = "G"
			willReset := true
			if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
				lastKeyPress = "Shift + G"
				willReset = false
			}
			go mapGenerationRoutine(willReset)
			lastKeyPressTime = time.Now()
		}

	}
	rl.CloseWindow()
}
