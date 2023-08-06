package osuautodeafen

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/micmonay/keybd_event"
)

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		if f.Name == "static" {
			if _, err := os.Stat("./deps/static"); errors.Is(err, os.ErrNotExist) {
				err := extractAndWriteFile(f)
				if err != nil {
					return err
				}
			}
		} else {
			err := extractAndWriteFile(f)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func GenerateKeybonding(deafenKey string) keybd_event.KeyBonding {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	if deafenKey == "" {
		fmt.Printf("[!] Deafen keybind is not set.")
		kb.SetKeys(keybd_event.VK_D)
		kb.HasALT(true)
		return kb
	}
	if strings.Contains(deafenKey, "ctrl+") {
		kb.HasCTRL(true)
		deafenKey = strings.ReplaceAll(deafenKey, "ctrl+", "")
	}
	if strings.Contains(deafenKey, "shift+") {
		kb.HasSHIFT(true)
		deafenKey = strings.ReplaceAll(deafenKey, "shift+", "")
	}
	if strings.Contains(deafenKey, "alt+") {
		kb.HasALT(true)
		deafenKey = strings.ReplaceAll(deafenKey, "alt+", "")
	}

	if strings.Contains(deafenKey, "a") {
		kb.SetKeys(keybd_event.VK_A)
	}
	if strings.Contains(deafenKey, "b") {
		kb.SetKeys(keybd_event.VK_B)
	}

	if strings.Contains(deafenKey, "c") {
		kb.SetKeys(keybd_event.VK_C)
	}

	if strings.Contains(deafenKey, "d") {
		kb.SetKeys(keybd_event.VK_D)
	}

	if strings.Contains(deafenKey, "e") {
		kb.SetKeys(keybd_event.VK_E)
	}

	if strings.Contains(deafenKey, "f") {
		kb.SetKeys(keybd_event.VK_F)
	}

	if strings.Contains(deafenKey, "g") {
		kb.SetKeys(keybd_event.VK_G)
	}

	if strings.Contains(deafenKey, "h") {
		kb.SetKeys(keybd_event.VK_H)
	}

	if strings.Contains(deafenKey, "i") {
		kb.SetKeys(keybd_event.VK_I)
	}

	if strings.Contains(deafenKey, "j") {
		kb.SetKeys(keybd_event.VK_J)
	}

	if strings.Contains(deafenKey, "k") {
		kb.SetKeys(keybd_event.VK_K)
	}

	if strings.Contains(deafenKey, "l") {
		kb.SetKeys(keybd_event.VK_L)
	}

	if strings.Contains(deafenKey, "m") {
		kb.SetKeys(keybd_event.VK_M)
	}

	if strings.Contains(deafenKey, "n") {
		kb.SetKeys(keybd_event.VK_N)
	}

	if strings.Contains(deafenKey, "o") {
		kb.SetKeys(keybd_event.VK_O)
	}

	if strings.Contains(deafenKey, "p") {
		kb.SetKeys(keybd_event.VK_P)
	}

	if strings.Contains(deafenKey, "q") {
		kb.SetKeys(keybd_event.VK_Q)
	}

	if strings.Contains(deafenKey, "r") {
		kb.SetKeys(keybd_event.VK_R)
	}

	if strings.Contains(deafenKey, "s") {
		kb.SetKeys(keybd_event.VK_S)
	}

	if strings.Contains(deafenKey, "t") {
		kb.SetKeys(keybd_event.VK_T)
	}

	if strings.Contains(deafenKey, "u") {
		kb.SetKeys(keybd_event.VK_U)
	}

	if strings.Contains(deafenKey, "v") {
		kb.SetKeys(keybd_event.VK_V)
	}

	if strings.Contains(deafenKey, "w") {
		kb.SetKeys(keybd_event.VK_W)
	}

	if strings.Contains(deafenKey, "x") {
		kb.SetKeys(keybd_event.VK_X)
	}

	if strings.Contains(deafenKey, "y") {
		kb.SetKeys(keybd_event.VK_Y)
	}

	if strings.Contains(deafenKey, "z") {
		kb.SetKeys(keybd_event.VK_Z)
	}

	if strings.Contains(deafenKey, "0") {
		kb.SetKeys(keybd_event.VK_0)
	}

	if strings.Contains(deafenKey, "1") {
		kb.SetKeys(keybd_event.VK_1)
	}

	if strings.Contains(deafenKey, "2") {
		kb.SetKeys(keybd_event.VK_2)
	}

	if strings.Contains(deafenKey, "3") {
		kb.SetKeys(keybd_event.VK_3)
	}

	if strings.Contains(deafenKey, "4") {
		kb.SetKeys(keybd_event.VK_4)
	}

	if strings.Contains(deafenKey, "5") {
		kb.SetKeys(keybd_event.VK_5)
	}

	if strings.Contains(deafenKey, "6") {
		kb.SetKeys(keybd_event.VK_6)
	}

	if strings.Contains(deafenKey, "7") {
		kb.SetKeys(keybd_event.VK_7)
	}

	if strings.Contains(deafenKey, "8") {
		kb.SetKeys(keybd_event.VK_8)
	}

	if strings.Contains(deafenKey, "9") {
		kb.SetKeys(keybd_event.VK_9)
	}

	if strings.Contains(deafenKey, "f1") {
		kb.SetKeys(keybd_event.VK_F1)
	}

	if strings.Contains(deafenKey, "f2") {
		kb.SetKeys(keybd_event.VK_F2)
	}

	if strings.Contains(deafenKey, "f3") {
		kb.SetKeys(keybd_event.VK_F3)
	}

	if strings.Contains(deafenKey, "f4") {
		kb.SetKeys(keybd_event.VK_F4)
	}

	if strings.Contains(deafenKey, "f5") {
		kb.SetKeys(keybd_event.VK_F5)
	}

	if strings.Contains(deafenKey, "f6") {
		kb.SetKeys(keybd_event.VK_F6)
	}

	if strings.Contains(deafenKey, "f7") {
		kb.SetKeys(keybd_event.VK_F7)
	}

	if strings.Contains(deafenKey, "f8") {
		kb.SetKeys(keybd_event.VK_F8)
	}

	if strings.Contains(deafenKey, "f9") {
		kb.SetKeys(keybd_event.VK_F9)
	}

	if strings.Contains(deafenKey, "f10") {
		kb.SetKeys(keybd_event.VK_F10)
	}

	if strings.Contains(deafenKey, "f11") {
		kb.SetKeys(keybd_event.VK_F11)
	}

	if strings.Contains(deafenKey, "f12") {
		kb.SetKeys(keybd_event.VK_F12)
	}

	if strings.Contains(deafenKey, "space") {
		kb.SetKeys(keybd_event.VK_SPACE)
	}

	if strings.Contains(deafenKey, "enter") {
		kb.SetKeys(keybd_event.VK_ENTER)
	}

	if strings.Contains(deafenKey, "backspace") {
		kb.SetKeys(keybd_event.VK_BACKSPACE)
	}

	if strings.Contains(deafenKey, "tab") {
		kb.SetKeys(keybd_event.VK_TAB)
	}

	if strings.Contains(deafenKey, "capslock") {
		kb.SetKeys(keybd_event.VK_CAPSLOCK)
	}

	if strings.Contains(deafenKey, "esc") {
		kb.SetKeys(keybd_event.VK_ESC)
	}

	if strings.Contains(deafenKey, "home") {
		kb.SetKeys(keybd_event.VK_HOME)
	}

	if strings.Contains(deafenKey, "end") {
		kb.SetKeys(keybd_event.VK_END)
	}

	if strings.Contains(deafenKey, "insert") {
		kb.SetKeys(keybd_event.VK_INSERT)
	}

	if strings.Contains(deafenKey, "del") {
		kb.SetKeys(keybd_event.VK_DELETE)
	}

	if strings.Contains(deafenKey, "pageup") {
		kb.SetKeys(keybd_event.VK_PAGEUP)
	}

	if strings.Contains(deafenKey, "pagedown") {
		kb.SetKeys(keybd_event.VK_PAGEDOWN)
	}

	if strings.Contains(deafenKey, "up") {
		kb.SetKeys(keybd_event.VK_UP)
	}

	if strings.Contains(deafenKey, "down") {
		kb.SetKeys(keybd_event.VK_DOWN)
	}

	if strings.Contains(deafenKey, "left") {
		kb.SetKeys(keybd_event.VK_LEFT)
	}

	if strings.Contains(deafenKey, "right") {
		kb.SetKeys(keybd_event.VK_RIGHT)
	}

	if strings.Contains(deafenKey, "arrowup") {
		kb.SetKeys(keybd_event.VK_UP)
	}

	if strings.Contains(deafenKey, "printscreen") {
		kb.SetKeys(keybd_event.VK_PRINT)
	}

	if strings.Contains(deafenKey, "scrolllock") {
		kb.SetKeys(keybd_event.VK_SCROLLLOCK)
	}

	if strings.Contains(deafenKey, "pause") {
		kb.SetKeys(keybd_event.VK_PAUSE)
	}

	if strings.Contains(deafenKey, "numlock") {
		kb.SetKeys(keybd_event.VK_NUMLOCK)
	}

	return kb
}
