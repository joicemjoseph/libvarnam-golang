package libvarnam-worker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/varnamproject/libvarnam-golang"
)

type varnamWorker struct {
	LanguageCode string
	Text         string
	ResultWriter io.Writer
	Done         chan struct{}
}

func NewVarnamWorker(langCode string, word string, writer io.Writer) varnamWorker {
	return varnamWorker{langCode, word, writer, make(chan struct{})}
}

var transliterationChannelMap = make(map[string]chan *varnamWorker)
var reverseTransliterationChannelMap = make(map[string]chan *varnamWorker)
var learningChannelMap = make(map[string]chan *varnamWorker)

func init() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	goRoutineCount := make(map[string]int)
	err := decoder.Decode(&goRoutineCount)
	if err != nil {
		fmt.Println("error:", err)
	}

	for language, count := range goRoutineCount {
		transliterationChannelMap[language] = make(chan *varnamWorker)
		reverseTransliterationChannelMap[language] = make(chan *varnamWorker)
		learningChannelMap[language] = make(chan *varnamWorker)
		go learningWorker(language, learningChannelMap[language])

		for i := 0; i < count; i++ {
			go transliterationWorker(language, transliterationChannelMap[language])
			go reverseTransliterationWorker(language, reverseTransliterationChannelMap[language])
		}
	}
}

func (v *varnamWorker) Transliterate() {
	transliterationChannelMap[v.LanguageCode] <- v
}

func (v *varnamWorker) ReverseTransliterate() {
	reverseTransliterationChannelMap[v.LanguageCode] <- v
}

func (v *varnamWorker) Learn() {
	learningChannelMap[v.LanguageCode] <- v
}

func transliterationWorker(langCode string, jobs <-chan *varnamWorker) {
	for j := range jobs {
		handle, error := libvarnam.Init(langCode)
		if error != nil {
			panic("Varnam not initialized")
		}
		outputs, varnamError := handle.Transliterate(j.Text)
		if varnamError != nil {
			panic("Transliteration failed for " + j.Text + " with error: " + varnamError.Error())
		}
		if j.ResultWriter == nil {
			panic("Writer is null")
		}
		fmt.Fprintf(j.ResultWriter, strings.Join(outputs, "\n"))
		close(j.Done)
	}
}

func reverseTransliterationWorker(langCode string, jobs <-chan *varnamWorker) {
	for j := range jobs {
		handle, error := libvarnam.Init(langCode)
		if error != nil {
			panic("Varnam not initialized")
		}
		output, varnamError := handle.ReverseTransliterate(j.Text)
		if varnamError != nil {
			panic("Reverse Transliteration failed for " + j.Text + " with error: " + varnamError.Error())
		}
		if j.ResultWriter == nil {
			panic("Writer is null")
		}
		fmt.Fprintf(j.ResultWriter, output)
		close(j.Done)
	}
}

func learningWorker(langCode string, jobs <-chan *varnamWorker) {
	var mutex = &sync.Mutex{}
	for j := range jobs {
		handle, error := libvarnam.Init(langCode)
		if error != nil {
			panic("Varnam not initialized")
		}
		mutex.Lock()
		varnamError := handle.Learn(j.Text)
		mutex.Unlock()
		if varnamError != nil {
			panic("Learning of " + j.Text + " failed with error: " + varnamError.Error())
		}
		if j.ResultWriter == nil {
			panic("Writer is null")
		}
		fmt.Fprintf(j.ResultWriter, "Done")
		close(j.Done)
	}
}