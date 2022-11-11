package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	utils "github.com/Nat3z/osudeafen/utils"
	"github.com/gorilla/websocket"
	"github.com/micmonay/keybd_event"
	"gopkg.in/ini.v1"
)

type ComboGosu struct {
	Current float64 `json:"current"`
	Max     float64 `json:"max"`
}

type BeatmapStatsGosu struct {
	MaxCombo float64 `json:"maxcombo"`
}
type BeatmapGosu struct {
	Stats BeatmapStatsGosu `json:"stats"`
	ID    int              `json:"id"`
	Time  TimeGosu         `json:"time"`
}

type TimeGosu struct {
	Current  float32 `json:"current"`
	Full     float32 `json:"full"`
	FirstObj float32 `json:"firstObj"`
}
type MenuGosu struct {
	BM    BeatmapGosu `json:"bm"`
	State int         `json:"state"`
}

type GameplayHitsGosu struct {
	Misses      float64 `json:"0"`
	Meh         float64 `json:"50"`
	Okay        float64 `json:"100"`
	Great       float64 `json:"300"`
	SliderBreak float64 `json:"sliderBreaks"`
}
type GameplayGosu struct {
	Name     string           `json:"name"`
	GameMode int              `json:"gamemode"`
	Score    float64          `json:"score"`
	Combo    ComboGosu        `json:"combo"`
	Accuracy float64          `json:"accuracy"`
	Hits     GameplayHitsGosu `json:"hits"`
}

type GoSuMemory struct {
	Gameplay GameplayGosu `json:"gameplay"`
	Menu     MenuGosu     `json:"menu"`
	Error    string       `json:"error"`
}

type GeneralSettings struct {
	Name                         string `ini:"username"`
	StartGosuMemoryAutomatically bool   `ini:"startgosumemory"`
}

type GameplaySettings struct {
	DeafenPercent       float64 `ini:"deafenpercent"`
	UndeafenAfterMisses float64 `ini:"undeafenmiss"`
}
type Settings struct {
	Gameplay GameplaySettings `ini:"gameplay"`
	General  GeneralSettings  `ini:"general"`
}

var addr = "localhost:24050"
var alreadyDeafened = false

var state int = 0
var recentlyjoined = false
var alreadyDetectedRestart = false
var inbeatmap = false
var misses float64 = 0

// true for deafen
// false for undeafen
func deafenOrUndeafen(kb keybd_event.KeyBonding, expect bool) {

	if alreadyDeafened {
		// if expecting a deafen, dont do anything.
		if expect {
			return
		}
		fmt.Println("| [KP] UNDEAFEN")
		kb.Launching()
	} else {
		// if expecting an undeafen, dont do anything.
		if !expect {
			return
		}
		fmt.Println("| [KP] DEAFEN")
		kb.Launching()
	}

	alreadyDeafened = !alreadyDeafened
}

func loadConfig() Settings {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("[!!] No config.ini found! Creating a config.ini...")
		out, _ := os.Create("config.ini")
		resp, err := http.Get("https://raw.githubusercontent.com/nat3z/osuautodeafen/master/config.ini.temp")
		if err != nil {
			fmt.Println("[!!] Unable to get template for osuautodeafen. Please connect to the internet and try again later.")
			os.Exit(1)
		}
		// tempout, _ := os.ReadFile("config.ini.temp")
		// temp := string(tempout)
		temp, _ := io.ReadAll(resp.Body)
		out.Write(([]byte)(temp))
		out.Close()
		fmt.Println("[#] Config.ini has been created! Please setup the config file and launch osuautodeafen.")
		time.Sleep(5 * time.Second)
		os.Exit(0)
		return Settings{}
	}
	var settings = new(Settings)
	cfg.MapTo(&settings)

	return *settings
}

func shutdown(cmnd exec.Cmd) {
	if err := cmnd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
	os.Exit(0)
}

var timesincelastws int64 = 0

func main() {
	fmt.Printf("[#] Checking for Updates...\n")
	utils.CheckVersion()
	utils.CheckVersionGosu()
	var config = loadConfig()

	// if start gosumemory automatically is on, then start process
	cmnd := exec.Command("./deps/gosumemory.exe")
	if config.General.StartGosuMemoryAutomatically {
		fmt.Printf("[#] Starting GosuMemory... \n")
		cmnd.Start()
		time.Sleep(2 * time.Second)
	}

	deafenKeybind := "alt+d"
	kb, err := keybd_event.NewKeyBonding()

	if err != nil {
		panic(err)
	}
	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_D)
	kb.HasALT(true)

	fmt.Printf("[!] Deafen keybind will be %s. Please make sure that your deafen keybind is set to this.\n", deafenKeybind)

	urlParsed := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	ws, _, err := websocket.DefaultDialer.Dial(urlParsed.String(), nil)

	if err != nil {
		fmt.Println("[!!] Error when connecting to GosuMemory. Please make sure that GosuMemory is open and is connected to osu!")
		shutdown(*cmnd)
		return
	}
	fmt.Println("[!] Connected to GosuMemory. Make sure that it stays on when playing osu!")
	fmt.Println("[!] Playing as", config.General.Name)

	timesincelastws = time.Now().Unix()

	go func() {
		for {
			if time.Now().Unix()-timesincelastws > 1 {
				fmt.Println("[!!] osu! has closed. Now stopping osu! Auto Deafen...")
				shutdown(*cmnd)
				break
			}
		}
	}()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("[!!] Error reading: ", err)
			break
		}
		var gosuResponse GoSuMemory
		jsonerr := json.Unmarshal(message, &gosuResponse)
		if jsonerr != nil {
			fmt.Println("[!!] ", jsonerr)
		} else {

			timesincelastws = time.Now().Unix()

			if gosuResponse.Gameplay.Name == config.General.Name && inbeatmap {

				if gosuResponse.Menu.BM.Time.Current > 1 && (recentlyjoined || alreadyDetectedRestart) {
					recentlyjoined = false
					alreadyDetectedRestart = false
				}

				if gosuResponse.Gameplay.Hits.Misses-misses != 0 {
					fmt.Println("| Missed, Broke, or lost combo. Incrementing miss count.")
					misses = gosuResponse.Gameplay.Hits.Misses
				}

				if misses >= config.Gameplay.UndeafenAfterMisses && alreadyDeafened {
					fmt.Printf("| Missed too many times (%sx) for undeafen. Now undeafening..\n", fmt.Sprint(config.Gameplay.UndeafenAfterMisses))
					deafenOrUndeafen(kb, false)
				}

				if gosuResponse.Gameplay.Score == 0 && gosuResponse.Gameplay.Accuracy == 0 && gosuResponse.Gameplay.Combo.Current == 0 && !recentlyjoined && !alreadyDetectedRestart {
					fmt.Println("| Detected that the user has restarted map. Attempting to undeafen..")
					misses = 0
					alreadyDetectedRestart = true
					deafenOrUndeafen(kb, false)
				} else if math.Floor(gosuResponse.Menu.BM.Stats.MaxCombo*config.Gameplay.DeafenPercent) < gosuResponse.Gameplay.Combo.Current && !alreadyDeafened && inbeatmap && misses == 0 {
					fmt.Println("| Reached max combo treshold for map. Now deafening..")
					deafenOrUndeafen(kb, true)
				}
			}

			if gosuResponse.Menu.State == 2 && state != 2 {
				fmt.Println("[#] Detected Beatmap Join")
				inbeatmap = true
				recentlyjoined = true
			} else if state == 2 && gosuResponse.Menu.State != 2 && inbeatmap {
				fmt.Println("[#] Detected Beatmap Exit")
				inbeatmap = false
				misses = 0
				deafenOrUndeafen(kb, false)
			}
			state = gosuResponse.Menu.State
		}
	}
	shutdown(*cmnd)
}
