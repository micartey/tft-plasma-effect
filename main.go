package main

import (
	"image/color"
	"log"
	"machine"

	"tinygo.org/x/drivers/gc9a01"
)

const (
	RESETPIN = machine.GPIO12
	CSPIN    = machine.GPIO9
	DCPIN    = machine.GPIO8
	BLPIN    = machine.GPIO25
)

func main() {
	spi := machine.SPI1
	conf := machine.SPIConfig{
		Frequency: 40 * machine.MHz,
	}
	if err := spi.Configure(conf); err != nil {
		log.Fatal(err)
	}

	lcd := gc9a01.New(spi, RESETPIN, DCPIN, CSPIN, BLPIN)
	lcd.Configure(gc9a01.Config{})
	lcd.IsBGR(true)

	width, height := lcd.Size()

	cx, cy := float32(width)/2, float32(height)/2
	// Corrected radius: A 240px wide display has a radius of 120px.
	radius := float32(width)

	lcd.FillRectangle(0, 0, width, height, color.RGBA{0, 0, 0, 255})
	lcd.FillRectangle(width/2-10, height/2-10, 20, 20, color.RGBA{255, 255, 255, 255})

	var t float32 = 1

	for {
		t += 0.05

		generatePlasma(lcd, width, height, cx, cy, radius, t, 8)

		log.Println("Frame drawn")
	}
}

func generatePlasma(lcd gc9a01.Device, width, height int16, cx, cy, radius, t float32, resolution int16) {
	// Define the size of the dark core (0.0 = center, 0.5 = 50% of radius)
	// You can adjust this value.
	const coreSize float32 = -1

	// Pre-calculate 1.0 / (1.0 - coreSize) for the remapping calculation
	// This avoids division inside the loop.
	var remapFactor float32 = 1.0
	if 1.0-coreSize > 0.001 { // Avoid division by zero if coreSize is 1.0
		remapFactor = 1.0 / (1.0 - coreSize)
	}

	for y := int16(0); y < height; y += resolution {
		for x := int16(0); x < width; x += resolution {
			dx := float32(x) - cx
			dy := float32(y) - cy
			dist := sqrt(dx*dx + dy*dy)

			var c color.RGBA

			if dist > radius || dist < radius*coreSize {
				c = color.RGBA{0, 0, 0, 255}
			} else {
				// True plasma calculation
				const angle float32 = /* atan2(dy, dx) */ 1
				radiusNorm := dist / radius // 0 at center, 1 at edge

				// Combine multiple waves
				v := sin(5*radiusNorm - t)
				v += sin(3*angle + t*1.5)
				v += sin(4*(dx+dy)/float32(width) + t*0.7)

				// Normalize to 0-1
				vn := (v + 3) / 6

				// --- MODIFIED SECTION ---

				// Remap radiusNorm to account for the core size.
				// 1. Subtract coreSize: (e.g., 0.0 -> -0.3,  0.3 -> 0.0,  1.0 -> 0.7)
				// 2. Multiply by remapFactor: (e.g., -0.3 -> <0,  0.0 -> 0.0,  0.7 -> 1.0)
				remappedNorm := (radiusNorm - coreSize) * remapFactor

				// Clamp the result (handles the core area, maps all <0 to 0)
				remappedNorm = clamp(remappedNorm, 0, 1)

				// Use sqrt() on the remapped value for the vibrant curve
				fade := sqrt(remappedNorm) * 2.0
				// --- END MODIFIED SECTION ---

				// Map to rainbow
				r := uint8(clamp((sin(vn*pi*2)+1)/2*255*fade, 0, 255))
				g := uint8(clamp((sin(vn*pi*2+2*pi/3)+1)/2*255*fade, 0, 255))
				b := uint8(clamp((sin(vn*pi*2+4*pi/3)+1)/2*255*fade, 0, 255))

				c = color.RGBA{r, g, b, 255}
			}

			lcd.FillRectangle(x, y, resolution, resolution, c)
		}
	}
}
