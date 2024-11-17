package rayfish

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/rwxrob/bonzai/anim"
)

// getTerminalSize retrieves current terminal dimensions
func getTerminalSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// fallback to some reasonable defaults if we can't get the size
		return 80, 24
	}
	return width, height
}

// clamp ensures a value stays within given range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// colorToANSI converts Color to ANSI terminal color string
func colorToANSI(c Color) string {
	r := int(c.R * 255)
	g := int(c.G * 255)
	b := int(c.B * 255)
	r = clamp(r, 0, 255)
	g = clamp(g, 0, 255)
	b = clamp(b, 0, 255)

	if r == 0 && g == 0 && b == 0 {
		return " " // Just space for black (eyes)
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s", r, g, b, "X")
}

func RunAnimation(numFish int, ground bool) {
	rand.Seed(time.Now().UnixNano())

	err := anim.SimpleAnimationScreen()
	if err != nil {
		fmt.Printf("Error initializing animation screen: %v\n", err)
		return
	}

	var frameBuilder strings.Builder
	frameBuilder.Grow(16384)

	scene, light := CreateScene(numFish, ground)

	lastTime := time.Now()
	for {
		currentTime := time.Now()
		deltaTime := float64(currentTime.Sub(lastTime).Seconds())
		lastTime = currentTime

		width, height := getTerminalSize()

		smallestDimension := width
		if height*2 < smallestDimension {
			smallestDimension = height * 2
		}

		aspectRatio := float64(width) / float64(height*2)

		frameBuilder.Reset()
		frameBuilder.WriteString("\033[H")
		frameBuilder.WriteString("\033[0m")

		UpdateScene(scene, deltaTime)

		line := make([]string, width+1)
		line[width] = "\n"

		lineIndex := 0
		RenderScene(scene, light, width, height, aspectRatio, func(x, y int, color Color) {
			if x == 0 && y > 0 {
				for _, s := range line {
					frameBuilder.WriteString(s)
				}
				lineIndex = 0
			}
			line[lineIndex] = colorToANSI(color)
			lineIndex++
		})

		for _, s := range line {
			frameBuilder.WriteString(s)
		}

		fmt.Print(frameBuilder.String())

		time.Sleep(33 * time.Millisecond)
	}
}
