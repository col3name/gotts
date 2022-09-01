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
 *  speech :+Speech{Folder: "audio", Language: voices.EnglishUK, Volume: 0, Speed: 1}
 * or
 * 	speech := NewSpeech(voices.Japanese, 0)
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Handler  handlers.PlayerInterface
	Volume   float64
	Speed    float64
}

const AudioFolder = "audio"

func NewSpeech(language string, volume float64) *Speech {
	return &Speech{Folder: AudioFolder, Language: language, Volume: volume, Speed: 1}
}

func (speech *Speech) CreateSpeechFile(text string, fileName string) (string, error) {
	err := speech.createFolderIfNotExists(speech.Folder)
	if err != nil {
		return "", err
	}

	filePath := speech.Folder + "/" + fileName + ".mp3"
	err = speech.downloadIfNotExists(filePath, text)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

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

func (speech *Speech) downloadIfNotExists(fileName string, text string) error {
	f, err := os.Open(fileName)
	if err != nil {
		response, err := http.Get(speech.getTranslatedFileURL(text))
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

func (speech *Speech) getTranslatedFileURL(text string) string {
	return fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
}

func (speech *Speech) generateHashName(name string) string {
	hash := md5.Sum([]byte(name))
	return fmt.Sprintf("%s_%s", speech.Language, hex.EncodeToString(hash[:]))
}
