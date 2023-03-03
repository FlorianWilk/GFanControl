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

func onReady() {
	mimage := image.NewRGBA(image.Rect(0, 0, 16, 16))
	systray.SetIcon(draws(mimage, 0))
	systray.SetTitle("")
	systray.SetTooltip("")
	go func() {
		for {
			time.Sleep(5 * time.Second)
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
				speed = 48 //25
			} else {
				speed = 0
			}
			systray.SetIcon(draws(mimage, speed))

			log.Printf("%d degree. Speed: %d", temp, speed)
			setSpeed(0, speed)
			setSpeed(1, speed)
		}
	}()
}

func ncmd(cmd string) string {
	return fmt.Sprintf("%s %s", command, cmd)
}

func rncmd(cmd2 string) string {
	cmd := ncmd(cmd2)
	log.Println(cmd2)
	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Printf("error %s while executing command: %s", err, string(out))
		return ""
	}
	return string(out)
}

func getTemp() int {
	out := rncmd(`-q="[gpu:0]/GPUCoreTemp" -t -c ":1"`)
	sout := strings.TrimRight(out, "\n")
	temp, err := strconv.Atoi(sout)
	if err != nil {
		log.Println("unable to parse core temp")
	}
	return temp
}

func setFanControl(gpu int, val int) {
	rncmd(fmt.Sprintf(`-a="[gpu:%d]/GPUFanControlState=%d"`, gpu, val))
}

func setSpeed(fan int, speed int) {
	rncmd(fmt.Sprintf(`-a="[fan:%d]/GPUTargetFanSpeed=%d"`, fan, speed))
}

func draws(mimage *image.RGBA, val int) []byte {
	col2 := color.RGBA{00, 0xa0, 0xff, 255}
	col := color.RGBA{80, 80, 80, 50}
	bcol := color.RGBA{255, 255, 255, 0}
	draw.Draw(mimage, mimage.Bounds(), &image.Uniform{bcol}, image.Point{}, draw.Src)
	mm := 7
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {

			xd := float64(x - mm)
			yd := float64(y - mm)
			a := math.Atan2(xd, yd) * (180 / math.Pi)
			d := math.Sqrt(float64(xd*xd + yd*yd))
			// if d < 3.0 && d > 3 {
			// 	if a > -180 && a < 180-(360.0/100.0*float64(val)) {
			// 		//					mimage.Set(x, y, col)

			// 	} else {
			// 		//					mimage.Set(x, y, col)
			// 	}
			// } else
			if d >= 4 && d <= 6.0 {
				if a > -180 && a < 180-(360.0/100.0*float64(val)) {
					mimage.Set(x, y, col)

				} else {
					mimage.Set(x, y, col2)
				}
			} else {
				mimage.Set(x, y, bcol)
			}
		}
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, mimage)
	return buf.Bytes()
}

func onExit() {
}

func main() {
	systray.Run(onReady, onExit)
}
