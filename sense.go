package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"time"

	"github.com/jweissig/amg8833"
	"github.com/mstahl/tsl2591"
)

// Reading - struct for sensor readings
type Reading struct {
	Timestamp time.Time
	Pixels    []float64
	Lux       float64
}

func main() {

	// thermal sensor
	amg, err := amg8833.NewAMG8833(&amg8833.Opts{
		Device: "/dev/i2c-0",
		Mode:   amg8833.AMG88xxNormalMode,
		Reset:  amg8833.AMG88xxInitialReset,
		FPS:    amg8833.AMG88xxFPS1,
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

	// start http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// lux
		channel0, channel1 := tsl.GetFullLuminosity()
		lux := tsl.CalculateLux(channel0, channel1)

		// check for NaN reading (when little to no light)
		if math.IsNaN(lux) == true {
			lux = 0.0
		}

		// pixels
		grid := amg.ReadPixels()
		pixels := ""
		for i, p := range grid {
			pixels = pixels + fmt.Sprintf("%3.2f ", p)
			// line break into 8x8 grid
			if i == 7 || i == 15 || i == 23 || i == 31 || i == 39 || i == 47 || i == 55 || i == 63 {
				pixels = pixels + fmt.Sprintf("\n")
			}
		}

		data := struct {
			Timestamp time.Time
			Pixels    string
			Lux       float64
		}{
			Timestamp: time.Now(),
			Pixels:    pixels,
			Lux:       lux,
		}

		t, _ := template.ParseFiles("html/index.html")
		t.Execute(w, data)
	})

	fs := http.FileServer(http.Dir("html/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	go http.ListenAndServe(":8787", nil)

	// how often to record sensor data
	ticker := time.NewTicker(1 * time.Minute)

	// loop
	for {

		// store sensor data
		var reading Reading

		// set timestamp
		reading.Timestamp = time.Now()

		// populate pixels
		reading.Pixels = amg.ReadPixels()

		// populate lux
		channel0, channel1 := tsl.GetFullLuminosity()
		//fmt.Println(tsl.CalculateLux(channel0, channel1))
		reading.Lux = tsl.CalculateLux(channel0, channel1)

		// check for NaN reading (when little to no light)
		if math.IsNaN(reading.Lux) == true {
			reading.Lux = 0.0
		}

		// encode as json
		data, err := json.Marshal(reading)
		if err != nil {
			fmt.Println(err)
		}

		// print to console
		fmt.Println(string(data))

		<-ticker.C
	}

}
