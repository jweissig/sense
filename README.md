Prototype data logger using a Raspberry Pi Zero W, AMG8833 8x8 Thermal Camera Sensor, and a TSL2591 High Dynamic Range Digital Light Sensor. The idea is to put one of these in a room and detect visible and non-visible light. This problem will allow you to connect to a web server on the Raspberry Pi and view the current status. But, the program will also log event each minute in JSON format to the console for later analysis.

## Setting up

You'll need to have both the AMG8833 and TSL2591 sensors connected to the I2C bus on the Raspberry Pi. I am using two buses since I just connected the sensors directly and didn't want to use an extra break out. You could easily just use a single I2C bus here. The drivers allow you to pick the bus though. Use `i2cdetect -y 0` or `i2cdetect -y 1` to verify you have both sensors connected.

    root@raspberrypi:~/sense# i2cdetect -y 0
         0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
    00:          -- -- -- -- -- -- -- -- -- -- -- -- --
    10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    60: -- -- -- -- -- -- -- -- -- 69 -- -- -- -- -- --
    70: -- -- -- -- -- -- -- --

    root@raspberrypi:~/sense# i2cdetect -y 1
         0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
    00:          -- -- -- -- -- -- -- -- -- -- -- -- --
    10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    20: -- -- -- -- -- -- -- -- -- 29 -- -- -- -- -- --
    30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
    70: -- -- -- -- -- -- -- --

## Installation

    go get -u github.com/jweissig/amg8833
    go get -u github.com/mstahl/tsl2591

## Usage

You can connect to a web server on localhost:8787 that will show you the Lux and Thermal sensor data in a gnuplot window. This is pretty useful for hacking around. Data will also be logged to the console which you can pipe/redirect into a logging file.

    go run sense.go
