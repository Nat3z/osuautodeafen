package osuautodeafen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/jxeng/shortcut"
	"github.com/ncruces/zenity"
)

type Message struct {
	Type  string `json:"message"`
	Value string `json:"value"`
}

type GeneralSettings struct {
	Name                         string `json:"username"`
	StartGosuMemoryAutomatically bool   `json:"startgosumemory"`
	DeafenKey                    string `json:"deafenkey"`
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

var resources = []string{
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/app.js",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/index.html",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/style.css",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/slider.css",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/assets/logo-not-transparent.png",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/version.txt",
	"https://raw.githubusercontent.com/Nat3z/osuautodeafen/master/resources/app/osu.ico",
}

func DownloadResources() {
	// download the resources
	fmt.Println("[#] Downloading resources..")
	// check if a resources folder exists
	if _, err := os.Stat("./resources/"); os.IsNotExist(err) {
		os.Mkdir("./resources", os.ModeAppend)
		os.Mkdir("./resources/app", os.ModeAppend)
	}
	for _, resource := range resources {
		resp, err := http.Get(resource)
		if err != nil {
			fmt.Println("[!!] Error occurred when downloading GosuMemory.")
			return
		}

		defer resp.Body.Close()
		bodyEncoded, _ := io.ReadAll(resp.Body)
		// create the file in the resources/app
		var fileName = resource[strings.LastIndex(resource, "/")+1:]
		out, err := os.Create("./resources/app/" + fileName)
		if err != nil {
			fmt.Println("[!!] Error occurred when creating file.")
			return
		}
		defer out.Close()
		// write the body to the file
		out.Write(bodyEncoded)
		fmt.Println("[#] Downloaded " + fileName)
	}
	fmt.Println("[#] Finished downloading resources.")
}

func CreateWindow(settings Settings, isFirstLoad bool) {

	if WindowAlreadyOpened {
		return
	}

	// see if the ./resources/app/index.html file exists
	_, err := os.Stat("./resources/app/index.html")
	if os.IsNotExist(err) {
		fmt.Println("[#] Preparing to download resources..")
		// download the resources
		DownloadResources()
	}

	// check if the file version.txt exists in resources/app
	_, err = os.Stat("./resources/app/version.txt")
	if os.IsNotExist(err) {
		fmt.Println("[!!] Error occurred when checking version.")
		return
	}
	// read the file
	versionFile, err := os.Open("./resources/app/version.txt")
	if err != nil {
		fmt.Println("[!!] Error occurred when checking version.")
		return
	}

	version, err := io.ReadAll(versionFile)
	if err != nil {
		fmt.Println("[!!] Error occurred when checking version.")
		return
	}
	// check if the version is the same as the one online
	resp, err := http.Get("https://raw.githubusercontent.com/Nat3z/osuautodeafen/future/resources/app/version.txt")
	if err != nil {
		fmt.Println("[!!] Error occurred when checking version.")
		return
	}
	defer resp.Body.Close()
	bodyEncoded, _ := io.ReadAll(resp.Body)
	if string(version) != string(bodyEncoded) {
		fmt.Println("[#] New version available, downloading..")
		// download the resources
		DownloadResources()
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
		Height:      astikit.IntPtr(205),
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
		var message SettingAsMessage
		json.Unmarshal([]byte(s), &message)
		if message.Type == "generate-shortcut" {
			// ask the user for where the "osu!.exe" file is
			// get the username of the user
			var username = os.Getenv("USERNAME")
			filedialog, diagerr := zenity.SelectFile(zenity.Title("Select osu!.exe"), zenity.FileFilters{
				zenity.FileFilter{
					Name: "osu!.exe",
				},
			}, zenity.Filename("C:\\"+username+"\\AppData\\Local\\osu!\\osu!.exe"))
			if diagerr != nil {
				fmt.Println("[!!] Error occurred when opening file dialog.")
				return "ERROR"
			}
			// get the path of the file
			var path = filedialog
			// get the local path of this .exe file
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}

			// get the path of the file (including the file name)
			exPath := filepath.Dir(ex) + "\\" + filepath.Base(ex)
			// create the shortcut
			//
			generatedShortcut := shortcut.Shortcut{
				ShortcutPath:     "C:\\Users\\" + username + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\osu! Auto Deafen.lnk",
				Target:           exPath,
				IconLocation:     filepath.Dir(ex) + "\\resources\\app\\osu.ico",
				Arguments:        "--open \"" + path + "\"",
				WorkingDirectory: filepath.Dir(ex),
			}
			shortcut.Create(generatedShortcut)
			return "SUCCESS"
		}
		remarshal, _ := json.Marshal(message.Value)
		if message.Value.General.DeafenKey == "" {
			message.Value.General.DeafenKey = "alt+d"
		}

		out, _ := os.Create("config.json")
		out.Write(([]byte)(remarshal))
		out.Close()

		return "SUCCESS"
	})
	// Blocking pattern
	a.Wait()
}
