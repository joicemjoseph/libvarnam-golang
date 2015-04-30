package libvarnam

import (
	"os"
	"testing"
)

func initVarnam(schemeId string, t *testing.T) *Varnam {
	varnam, err := Init(schemeId)
	if err != nil {
		t.Errorf("Expected init to run, but failing with: %s", err.Error())
	}
	return varnam
}

func TestInit(t *testing.T) {
	initVarnam("ml", t)
}

func TestInitWithIncorrectIdentifierCode(t *testing.T) {
	_, err := Init("ml-nonexisting")
	if err == nil {
		t.Error("Expected init to fail when lang code is incorrect")
	}
	expectedErrorMessage := "Failed to find symbols file for: ml-nonexisting"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message to be: %s, but was: %s", expectedErrorMessage, err.Error())
	}
}

func TestGetSuggestionsFilePath(t *testing.T) {
	varnam := initVarnam("ml", t)
	suggestionsFilePath := varnam.GetSuggestionsFilePath()
	if _, err := os.Stat(suggestionsFilePath); os.IsNotExist(err) {
		t.Errorf("%s: Suggestions file does not exists", suggestionsFilePath)
	}
}

func TestGetCorpusDetails(t *testing.T) {
	varnam := initVarnam("hi", t)
	_, err := varnam.GetCorpusDetails()
	if err != nil {
		t.Errorf("Failed to get corpus details. Got error: %s", err.Error())
	}
}

func TestTransliterate(t *testing.T) {
	varnam := initVarnam("ml", t)
	words, err := varnam.Transliterate("navaneeth")
	if err != nil {
		t.Errorf("Failed to perform transliteration. Got: %s", err.Error())
	}

	if len(words) == 0 {
		t.Errorf("Failed to perform transliteration. Got 0 words")
	}
}

func TestLearn(t *testing.T) {
	varnam := initVarnam("hi", t)
	err := varnam.Learn("भारत")
	if err != nil {
		t.Errorf("Failed to perform learn. Got error: %s", err.Error())
	}
}

func TestLearnFailure(t *testing.T) {
	varnam := initVarnam("hi", t)
	err := varnam.Learn("foo")
	if err == nil {
		t.Errorf("Expected learn to fail, but didn't fail")
	}
}

func TestGetAllSchemeDetails(t *testing.T) {
	scheme_details := GetAllSchemeDetails()
	if len(scheme_details) == 0 {
		t.Errorf("Expected GetAllSchemeDetails to return atlest one scheme details. But returned none")
	}
}

func TestDestroy(t *testing.T) {
	varnam := initVarnam("hi", t)
	varnam.Destroy()
}
