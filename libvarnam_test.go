package libvarnam

import (
	"testing"
)

func initVarnam(langCode string, t *testing.T) *Varnam {
	varnam, err := Init(langCode)
	if err != nil {
		t.Errorf("Expected init to run, but failing with: %s", err.Error())
	}
	return varnam
}

func TestInit(t *testing.T) {
	initVarnam("ml", t)
}

func TestInitWithIncorrectLangCode(t *testing.T) {
	_, err := Init("ml-nonexisting")
	if err == nil {
		t.Error("Expected init to fail when lang code is incorrect")
	}
	expectedErrorMessage := "Failed to find symbols file for: ml-nonexisting"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message to be: %s, but was: %s", expectedErrorMessage, err.Error())
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

func TestDestroy(t *testing.T) {
	varnam := initVarnam("hi", t)
	varnam.Destroy()
}
