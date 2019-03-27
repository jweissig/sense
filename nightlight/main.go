// trigger on 2+ frames (200ms vs 100ms)

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/boombuler/led"
	"github.com/jweissig/amg8833"
	"github.com/mstahl/tsl2591"
)

// Reading - struct for sensor readings
type Reading struct {
	Timestamp time.Time
	Pixels    []float64
	Lux       float64
}

var (
	state     string // on/off
	col       string
	lighttype string
	COLOR_OFF = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

func main() {

	// thermal sensor
	amg, err := amg8833.NewAMG8833(&amg8833.Opts{
		Device: "/dev/i2c-0",
		Mode:   amg8833.AMG88xxNormalMode,
		Reset:  amg8833.AMG88xxInitialReset,
		FPS:    amg8833.AMG88xxFPS10,
	})
	if err != nil {
		panic(err)
	}

	// lux sensor
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.TSL2591_GAIN_LOW,
		Timing: tsl2591.TSL2591_INTEGRATIONTIME_600MS,
	})
	if err != nil {
		panic(err)
	}

	for devInfo := range led.Devices() {
		// keep sure to find BlinkStick devices only
		if devInfo.GetType() != led.BlinkStick {
			continue
		}

		dev, err := devInfo.Open()

		// stop on error
		if err != nil {
			fmt.Println(err)
			break
		}
		defer dev.Close()

		dev.SetKeepActive(true)
		dev.SetColor(COLOR_OFF)

		// how often to check sensor data
		ticker := time.NewTicker(100 * time.Millisecond)

		// loop
		for {

			// turn off led
			dev.SetColor(COLOR_OFF)

			// populate pixels
			pixels := amg.ReadPixels()

			// get lux
			channel0, channel1 := tsl.GetFullLuminosity()
			lux := tsl.CalculateLux(channel0, channel1)

			// check for NaN reading (when little to no light)
			if math.IsNaN(lux) == true {
				lux = 0.0
			}

			//fmt.Println(pixels)
			//fmt.Println(lux)

			// 20.00  19.75  20.25  20.5   21.00  21.25  22.00  22.25 [7]
			// 20.25  20.00  20.00  21.00  20.5   21.75  22.5   22.00 [15]
			// 20.25  20.5   20.75  21.00  20.75  22.5   22.00  23.00 [23]
			// 20.00  20.5   20.75  20.25  20.5   21.25  22.25  23.75 [31]
			// 20.75  20.75  20.75  20.75  21.00  21.25  22.5   24.00 [39]
			// 20.5   20.00  20.5   20.5   21.00  21.25  23.00  24.25 [47]
			// 20.75  20.5   20.75  20.75  20.5   20.75  23.00  23.5  [55]
			// 20.00  20.75  20.25  20.75  20.5   20.5   20.75  20.75 [63]

			// if pixels[39] > 22 || pixels[47] > 22 {
			//if lux == 0.0 && pixels[39] > 23 {
			if lux == 0.0 && pixels[39] > 23.25 {

				log.Print("tripped: %s", pixels)

				//"darkred": color.RGBA{0x8b, 0x00, 0x00, 0xff},
				//"red":     color.RGBA{0xff, 0x00, 0x00, 0xff},
				//"white":   color.RGBA{0xff, 0xff, 0xff, 0xff},
				//"off":     color.RGBA{0x00, 0x00, 0x00, 0xff},
				dev.SetColor(color.RGBA{0xff, 0x00, 0x00, 0xff}) // red
				time.Sleep(150 * time.Second)
				dev.SetColor(COLOR_OFF)

			}

			<-ticker.C

		}

	}

}
