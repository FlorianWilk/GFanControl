package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

const command = "/usr/bin/nvidia-settings"

const FAN_COUNT = 2
const ICON_SIZE = 16
const UPDATE_INTERVAL_SEC = 5
const TEMP_MIN = 40
const TEMP_MAX = 90

var bgColor = color.RGBA{0, 0, 0, 0}
var fanSpeedColor = color.RGBA{00, 0xa0, 0xff, 255}
var circleColor = color.RGBA{80, 80, 80, 50}

func onReady() {
	mimage := image.NewRGBA(image.Rect(0, 0, ICON_SIZE, ICON_SIZE))
	systray.SetIcon(draws(mimage, 0, 0))
	systray.SetTitle("")
	systray.SetTooltip("")
	sysInfo := systray.AddMenuItem("", "")
	go func() {
		for {
			temp := getTemp()
			setFanControl(0, 1)
			var speed int
			if temp > 75 {
				speed = 85
			} else if temp > 65 {
				speed = 70
			} else if temp > 60 {
				speed = 65
			} else if temp > 55 {
				speed = 55
			} else if temp > 45 {
				speed = 52
			} else if temp > 35 {
				speed = 48
			} else {
				speed = 0
			}
			systray.SetIcon(draws(mimage, speed, temp))

			s := fmt.Sprintf("Temp: %d° / Speed: %d%%", temp, speed)
			sysInfo.SetTitle(s)

			for i := 0; i < FAN_COUNT; i++ {
				setSpeed(i, speed)
			}
			time.Sleep(UPDATE_INTERVAL_SEC * time.Second)
		}
	}()
}

func runcmd(cmd2 string) string {
	cmd := fmt.Sprintf("%s %s", command, cmd2)
	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Printf("error %s while executing command: %s", err, string(out))
		return ""
	}
	return string(out)
}

func getTemp() int {
	out := runcmd(`-q="[gpu:0]/GPUCoreTemp" -t`)
	sout := strings.TrimRight(out, "\n")
	temp, err := strconv.Atoi(sout)
	if err != nil {
		log.Println("unable to parse core temp")
	}
	return temp
}

func setFanControl(gpu int, val int) {
	runcmd(fmt.Sprintf(`-a="[gpu:%d]/GPUFanControlState=%d"`, gpu, val))
}

func setSpeed(fan int, speed int) {
	runcmd(fmt.Sprintf(`-a="[fan:%d]/GPUTargetFanSpeed=%d"`, fan, speed))
}

func draws(mimage *image.RGBA, val int, temp int) []byte {

	col3 := color.NRGBA{0xff, 0, 0xff, 255}

	draw.Draw(mimage, mimage.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	mm := (ICON_SIZE / 2) - 1

	t2 := math.Max(float64(temp-TEMP_MIN), 0)
	dd := float64(TEMP_MAX - TEMP_MIN)
	col3.A = uint8((255.0 / dd) * math.Min(float64(t2), dd))

	for y := 0; y < ICON_SIZE; y++ {
		for x := 0; x < ICON_SIZE; x++ {

			xd := float64(x - mm)
			yd := float64(y - mm)
			a := math.Atan2(xd, yd) * (180 / math.Pi)
			d := math.Sqrt(float64(xd*xd + yd*yd))

			if d < 3.0 && d >= 0 {
				mimage.Set(x, y, col3)
			} else if d >= 4 && d <= 6.0 {
				if a > -180 && a < 180.0-(360.0/100.0*float64(val)) {
					mimage.Set(x, y, circleColor)
				} else {
					mimage.Set(x, y, fanSpeedColor)
				}
			} else {
				mimage.Set(x, y, bgColor)
			}
		}
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, mimage)
	return buf.Bytes()
}

func onExit() {
	setFanControl(0, 0)
}

func main() {
	systray.Run(onReady, onExit)
}
