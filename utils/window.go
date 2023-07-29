package osuautodeafen

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

type Message struct {
	Type  string `json:"message"`
	Value string `json:"value"`
}

type GeneralSettings struct {
	Name                         string `json:"username"`
	StartGosuMemoryAutomatically bool   `json:"startgosumemory"`
}

type GameplaySettings struct {
	DeafenPercent           float64 `json:"deafenpercent"`
	UndeafenAfterMisses     float64 `json:"undeafenmiss"`
	CountSliderBreaksAsMiss bool    `json:"countsliderbreakmiss"`
}
type Settings struct {
	Gameplay GameplaySettings `json:"gameplay"`
	General  GeneralSettings  `json:"general"`
}
type SettingAsMessage struct {
	Type  string   `json:"type"`
	Value Settings `json:"value"`
}

var State int = 0
var WindowAlreadyOpened = false

func CreateWindow(settings Settings, isFirstLoad bool) {
	if WindowAlreadyOpened {
		return
	}
	var a, _ = astilectron.New(nil, astilectron.Options{
		AppName:            "osuautodeafen",
		AppIconDefaultPath: "./resources/icon.png",
		VersionAstilectron: "0.33.0",
		VersionElectron:    "4.0.1",
	})
	defer a.Close()

	// Start astilectron
	a.Start()
	// 220 width by default
	var w, _ = a.NewWindow("./resources/app/index.html", &astilectron.WindowOptions{
		Height:      astikit.IntPtr(200),
		Width:       astikit.IntPtr(195),
		AlwaysOnTop: astikit.BoolPtr(true),
		Transparent: astikit.BoolPtr(true),
		Frame:       astikit.BoolPtr(false),
		Resizable:   astikit.BoolPtr(false),
		X:           astikit.IntPtr(15),
		Y:           astikit.IntPtr(100),
	})
	WindowAlreadyOpened = true
	w.Create()
	go func() {
		var closed = false
		for {
			if State != 0 && !closed {
				fmt.Println("[#] Closing Window..")
				w.Close()
				closed = true
			} else if closed && State == 0 {
				WindowAlreadyOpened = false
				break
			}
		}
	}()
	var settingTypeName = "load"
	if isFirstLoad {
		settingTypeName += "-FIRSTLOAD"
	}
	var message SettingAsMessage = SettingAsMessage{Type: settingTypeName, Value: settings}
	loadsettingsout, _ := json.Marshal(message)

	w.SendMessage(string(loadsettingsout), func(m *astilectron.EventMessage) {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process message
		fmt.Printf("[#] %s\n", s)
	})
	w.OnMessage(func(m *astilectron.EventMessage) (v interface{}) {
		var s string
		m.Unmarshal(&s)
		fmt.Println(s)
		var message SettingAsMessage
		json.Unmarshal([]byte(s), &message)
		remarshal, _ := json.Marshal(message.Value)
		out, _ := os.Create("config.json")
		out.Write(([]byte)(remarshal))
		out.Close()

		return "SUCCESS"
	})
	// Blocking pattern
	a.Wait()
}
