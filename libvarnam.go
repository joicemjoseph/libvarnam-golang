package libvarnam

// #cgo pkg-config: varnam
// #include <varnam.h>
import "C"
import "fmt"

type Varnam struct {
	handle *C.varnam
}

type VarnamError struct {
	errorCode int
	message   string
}

type SchemeDetails struct {
	LangCode     string
	Identifier   string
	DisplayName  string
	Author       string
	CompiledDate string
	IsStable     bool
}

func (e *VarnamError) Error() string {
	return e.message
}

func (v *Varnam) Transliterate(text string) ([]string, error) {
	var va *C.varray
	rc := C.varnam_transliterate(v.handle, C.CString(text), &va)
	if rc != C.VARNAM_SUCCESS {
		errorCode := (int)(rc)
		return nil, &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	var i C.int
	var array []string
	for i = 0; i < C.varray_length(va); i++ {
		word := (*C.vword)(C.varray_get(va, i))
		array = append(array, C.GoString(word.text))
	}
	return array, nil
}

func (v *Varnam) ReverseTransliterate(text string) (string, error) {
	var output *C.char
	rc := C.varnam_reverse_transliterate(v.handle, C.CString(text), &output)
	if rc != C.VARNAM_SUCCESS {
		errorCode := (int)(rc)
		return "", &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	return C.GoString(output), nil
}

func Init(schemeIdentifier string) (*Varnam, error) {
	var v *C.varnam
	var msg *C.char
	rc := C.varnam_init_from_id(C.CString(schemeIdentifier), &v, &msg)
	if rc != C.VARNAM_SUCCESS {
		return nil, &VarnamError{errorCode: (int)(rc), message: C.GoString(msg)}
	}
	return &Varnam{handle: v}, nil
}

func GetAllSchemeDetails() []*SchemeDetails {
	allSchemeDetails := C.varnam_get_all_scheme_details()
	if allSchemeDetails == nil {
		return []*SchemeDetails{}
	}

	var schemeDetails []*SchemeDetails
	length := int(C.varray_length(allSchemeDetails))
	for i := 0; i < length; i++ {
		detail := (*C.vscheme_details)(C.varray_get(allSchemeDetails, C.int(i)))
		schemeDetails = append(schemeDetails, &SchemeDetails{
			LangCode: C.GoString(detail.langCode), Identifier: C.GoString(detail.identifier),
			DisplayName: C.GoString(detail.displayName), Author: C.GoString(detail.author),
			CompiledDate: C.GoString(detail.compiledDate), IsStable: detail.isStable > 0})

		C.varnam_destroy_scheme_details(detail)
	}
	C.varray_free(allSchemeDetails, nil)
	return schemeDetails
}

func (v *Varnam) Learn(text string) error {
	rc := C.varnam_learn(v.handle, C.CString(text))
	if rc != 0 {
		errorCode := (int)(rc)
		return &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	return nil
}

func (v *Varnam) getVarnamError(errorCode int) string {
	errormessage := C.varnam_get_last_error(v.handle)
	varnamErrorMsg := C.GoString(errormessage)
	return fmt.Sprintf("%d:%s", errorCode, varnamErrorMsg)
}

func (v *Varnam) Destroy() {
	C.varnam_destroy(v.handle)
}
