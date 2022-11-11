package osuautodeafen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var version = "v1.0"

type ReleaseAsset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
	Name               string `json:"name"`
}
type Releases struct {
	Assets  []ReleaseAsset `json:"assets"`
	TagName string         `json:"tag_name"`
}

func CheckVersionGosu() bool {
	resp, err := http.Get("https://api.github.com/repos/l3lackShark/gosumemory/releases/latest")
	if err != nil {
		fmt.Println("[!!] Error occurred when checking for updates")
		return false
	}
	defer resp.Body.Close()
	var releases Releases
	bodyEncoded, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bodyEncoded, &releases)
	var downloadLink = ""
	for _, release := range releases.Assets {
		if strings.Contains(release.Name, "gosumemory_windows_") && !strings.Contains(release.Name, "amd") {
			downloadLink = release.BrowserDownloadURL
			break
		}
	}

	var releaseVersion = ""
	if _, err := os.Stat("./deps/"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("./deps", os.ModeAppend)
		out, err := os.Create("./deps/version.txt")
		if err != nil {
			fmt.Println("[!!] Error occurred when creating version.txt")
			return false
		}
		releaseVersion = ""
		defer out.Close()
		// downloadGosuMemory(downloadLink)
	} else {
		version, _ := os.ReadFile("./deps/version.txt")
		releaseVersion = string(version)
	}

	if releaseVersion == releases.TagName {
		fmt.Printf("[#] Up-to-date with gosumemory repo.\n")
		return true
	}

	var depserror = os.Remove("./deps/gosumemory.exe")
	if depserror != nil {
		fmt.Println("[+] Preparing for clean install of gosumemory.exe")
	}
	DownloadGosuMemory(downloadLink)
	out, _ := os.Create("./deps/version.txt")
	out.WriteString(releases.TagName)
	defer out.Close()

	return true
}

func CheckVersion() bool {
	resp, err := http.Get("https://api.github.com/repos/Nat3z/osuautodeafen/releases/latest")
	if err != nil {
		fmt.Println("[!!] Error occurred when checking for updates.")
		return false
	}
	defer resp.Body.Close()
	var releases Releases
	bodyEncoded, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bodyEncoded, &releases)

	if version == releases.TagName {
		fmt.Printf("[#] Up-to-date with osu! Auto Deafen.\n")
	} else {
		fmt.Println("======================\n   osu! Auto Deafen\n   UPDATE AVAILABLE\n======================")
		fmt.Printf("VERSION %s: https://github.com/Nat3z/osuautodeafen/releases/latest\n", releases.TagName)
	}
	return true
}
func DownloadGosuMemory(filepath string) bool {
	out, err := os.Create("./deps/gosumemory.zip")
	if err != nil {
		fmt.Println("[!!] Error occurred when creating gosumemory.zip")
		return false
	}
	defer out.Close()
	resp, err := http.Get(filepath)
	if err != nil {
		fmt.Println("[!!] Error occurred when downloading GosuMemory.")
		return false
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)

	Unzip("./deps/gosumemory.zip", "./deps/")
	// removeerr := os.Remove("./deps/gosumemory.zip")
	// if removeerr != nil {
	// 	fmt.Println("[!!] Error occurred when deleting gosumemory.zip: ", removeerr)
	// }
	fmt.Printf("[+] Update completed.\n")
	return true
}
