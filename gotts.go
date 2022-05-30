package gotts

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/col3name/gotts/handlers"
	"io"
	"net/http"
	"net/url"
	"os"
)

/**
 *
 * Use:
 *
 * speech := htgotts.Speech{Folder: "audio", Language: "en"}
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Handler  handlers.PlayerInterface
	Volume   float64
	Speed    float64
}

// Creates a speech file with a given name
func (speech *Speech) CreateSpeechFile(text string, fileName string) (string, error) {
	var err error

	f := speech.Folder + "/" + fileName + ".mp3"
	if err = speech.createFolderIfNotExists(speech.Folder); err != nil {
		return "", err
	}

	if err = speech.downloadIfNotExists(f, text); err != nil {
		return "", err
	}

	return f, nil
}

// Plays an existent .mp3 file
func (speech *Speech) PlaySpeechFile(fileName string) error {
	if speech.Handler == nil {
		speech.Handler = &handlers.BeepPlayer{Volume: speech.Volume, Speed: 1}
	}
	err := speech.Handler.Play(fileName)
	if err != nil {
		return err
	}
	speech.deleteFile(fileName)
	return nil
}

func (speech *Speech) deleteFile(fileName string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(path + "\\" + fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func (speech *Speech) Speak(text string) error {
	var err error
	generatedHashName := speech.generateHashName(text)

	fileName, err := speech.CreateSpeechFile(text, generatedHashName)
	if err != nil {
		return err
	}

	return speech.PlaySpeechFile(fileName)
}

func (speech *Speech) createFolderIfNotExists(folder string) error {
	dir, err := os.Open(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0700)
	}

	dir.Close()
	return nil
}

/**
 * Download the voice file if does not exists.
 */
func (speech *Speech) downloadIfNotExists(fileName string, text string) error {
	f, err := os.Open(fileName)
	if err != nil {
		urlString := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
		fmt.Println(urlString)
		response, err := http.Get(urlString)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		output, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(output, response.Body)
		return err
	}

	defer f.Close()
	return nil
}

func (speech *Speech) generateHashName(name string) string {
	hash := md5.Sum([]byte(name))
	return fmt.Sprintf("%s_%s", speech.Language, hex.EncodeToString(hash[:]))
}
