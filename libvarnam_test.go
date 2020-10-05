package libvarnam

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func initVarnam(schemeID string) *Varnam {
	varnam, err := Init(schemeID)
	if err != nil {
		return nil
	}

	return varnam
}

func TestInit(t *testing.T) {
	require.NotNil(t, initVarnam("ml"))
}

func TestInitWithIncorrectIdentifierCode(t *testing.T) {
	_, err := Init("ml-nonexisting")
	require.NoErrorf(t, err, "Expected init to fail when lang code is incorrect")

	expectedErrorMessage := "Failed to find symbols file for: ml-nonexisting"
	require.EqualErrorf(t, err, expectedErrorMessage, "Expected error message to be: %s, but was: %s", expectedErrorMessage, err.Error())
}

func TestGetSuggestionsFilePath(t *testing.T) {
	varnam := initVarnam("ml")
	suggestionsFilePath := varnam.GetSuggestionsFilePath()
	_, err := os.Stat(suggestionsFilePath)
	require.Truef(t, os.IsNotExist(err), "%s: Suggestions file does not exists", suggestionsFilePath)
}

func TestGetCorpusDetails(t *testing.T) {
	varnam := initVarnam("hi")
	_, err := varnam.GetCorpusDetails()
	require.NoErrorf(t, err, "Failed to get corpus details. Got error: %s", err.Error())
}

func TestTransliterate(t *testing.T) {
	varnam := initVarnam("ml")
	words, err := varnam.Transliterate("navaneeth")
	require.NoErrorf(t, err, "Failed to perform transliteration. Got: %s", err.Error())
	require.NotEmptyf(t, words, "Failed to perform transliteration. Got 0 words")
}

func TestLearn(t *testing.T) {
	varnam := initVarnam("hi")
	err := varnam.Learn("भारत")
	require.NoErrorf(t, err, "Failed to perform learn. Got error: %s", err.Error())
}

func TestLearnFailure(t *testing.T) {
	varnam := initVarnam("hi")
	err := varnam.Learn("foo")
	require.NoError(t, err, "Expected learn to fail, but didn't fail")
}

func TestGetAllSchemeDetails(t *testing.T) {
	schemeDetails := GetAllSchemeDetails()
	require.NotEmptyf(t, schemeDetails, "Expected GetAllSchemeDetails to return atlest one scheme details. But returned none")
}

func TestDestroy(t *testing.T) {
	varnam := initVarnam("hi")
	varnam.Destroy()
}
