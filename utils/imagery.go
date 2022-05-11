package utils

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	color2 "image/color"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

var (
	BasePath = "./assets/1.18-item-icons/"
)

func MapIcons() {
	items, _ := ioutil.ReadDir(BasePath)
	for _, item := range items {
		core.ItemIcons = append(core.ItemIcons, item.Name())
	}
	sort.StringSlice.Sort(core.ItemIcons)
	logger.Info("grabbed and sorted item images!")
}

func GenerateAdvancement(icon string, advType string, advancement string) {
	var (
		title string
		color color2.RGBA
	)

	switch advType {
	case "CHALLENGE":
		title = "Challenge Complete!"
		color = color2.RGBA{R: 255, G: 85, B: 255, A: 255}
		break
	case "GOAL":
		title = "Goal Reached!"
		color = color2.RGBA{R: 255, G: 255, B: 85, A: 255}
		break
	case "ADVANCEMENT":
		title = "Advancement Made!"
		color = color2.RGBA{R: 255, G: 255, B: 85, A: 255}
		break
	}

	dc := gg.NewContext(400, 80)
	bg, err := gg.LoadImage("./assets/advancement.png")
	if err != nil {
		logger.Error("Failed to load Advancement image! Is it missing?\n", err)
	}
	ic, err := gg.LoadImage(icon)
	if err != nil {
		logger.Error("Failed to load Icon! Is it missing?\n", err)
	}
	fontPath := filepath.Join("./assets/fonts/minecraft_font.ttf")
	if err := dc.LoadFontFace(fontPath, 24); err != nil {
		logger.Error("Failed to load font!", err)
	}

	x1 := 24.0 + 48
	y1 := 32.0 + 4
	x2 := 24.0 + 48
	y2 := 64.0 - 4

	newIc := imaging.Resize(ic, 42, 42, imaging.Box)

	dc.DrawImage(bg, 0, 0)
	dc.DrawImage(newIc, 16, 16)

	dc.SetColor(color)
	dc.DrawString(title, x1, y1)
	dc.SetColor(color2.White)
	dc.DrawString(strings.Split(advancement, "/")[1], x2, y2)
	if err := dc.SavePNG("./assets/generated/advancement.png"); err != nil {
		logger.Error("Failed to save image!\n", err)
	}
}