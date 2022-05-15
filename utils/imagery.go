package utils

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	color2 "image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

func GenerateAdvancement(icon string, advType string, advancement string, player string) {
	_, err := os.Stat("./assets/generated")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./assets/generated", os.ModePerm)
	}
	_, err = os.Stat(fmt.Sprintf("./assets/generated/%v", player))
	if os.IsNotExist(err) {
		_ = os.Mkdir(fmt.Sprintf("./assets/generated/%v", player), os.ModePerm)
	}
	var (
		title string
		color color2.RGBA
		gray  = color2.RGBA{R: 170, G: 170, B: 170, A: 255}
		W     = 400
		H     = 110
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
	case "TASK":
		title = "Advancement Made!"
		color = color2.RGBA{R: 255, G: 255, B: 85, A: 255}
		break
	}

	dc := gg.NewContext(W, H)
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

	dc.DrawImage(bg, 0, 0)
	dc.DrawImage(ic, 16, 16)

	dc.SetColor(color)
	dc.DrawString(title, x1, y1)
	dc.SetColor(color2.White)
	dc.DrawString(advancement, x2, y2)

	if err := dc.LoadFontFace(fontPath, 20); err != nil {
		logger.Error("Failed to load font!", err)
	}

	completew, _ := dc.MeasureString("completed by:")
	_, playerh := dc.MeasureString(player)

	x4 := 24.0
	y4 := float64(H) - playerh - 8
	x5 := x4 + completew + 4
	y5 := y4
	dc.DrawString(player, x5, y5)
	dc.SetColor(gray)
	dc.DrawString("completed by", x4, y4)
	if err := dc.SavePNG(fmt.Sprintf("./assets/generated/%v/%v_advancement.png", player, player)); err != nil {
		logger.Error("Failed to save image!\n", err)
	}
}
