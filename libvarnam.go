package libvarnam

// #cgo pkg-config: varnam
// #include <stdio.h>
// #include <varnam.h>
import "C"

var errorMessagesMap = map[int]string{
	1: "VARNAM MISUSE",
	2: "VARNAM ERROR",
	3: "VARNAM ARGUMENTS ERROR",
	4: "VARNAM MEMORY ERROR",
	5: "VARNAM PARTIAL RENDERING",
	6: "VARNAM STORAGE ERROR",
	7: "VARNAM INVALID CONFIG",
	8: "VARNAM STEMRULE HIT",
	9: "VARNAM STEMRULE MISS",
}

type Varnam struct {
	handle *C.varnam
}

type VarnamError struct {
	errorCode int
	message   string
}

func (e *VarnamError) Error() string {
	return e.message
}

func (v *Varnam) Transliterate(text string) ([]string, *VarnamError) {
	var va *C.varray
	rc := C.varnam_transliterate(v.handle, C.CString(text), &va)
	if rc != 0 {
		errorCode := (int)(rc)
		return []string{}, &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	var i C.int
	var array []string
	for i = 0; i < C.varray_length(va); i++ {
		word := (*C.vword)(C.varray_get(va, i))
		array = append(array, C.GoString(word.text))
	}
	return array, nil
}

func (v *Varnam) ReverseTransliterate(text string) (string, *VarnamError) {
	var output *C.char
	rc := C.varnam_reverse_transliterate(v.handle, C.CString(text), &output)
	if rc != 0 {
		errorCode := (int)(rc)
		return "", &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	return C.GoString(output), nil
}

func Init(langCode string) (*Varnam, *VarnamError) {
	var v *C.varnam
	var msg *C.char
	rc := C.varnam_init_from_lang(C.CString(langCode), &v, &msg)
	if rc != 0 {
		return nil, &VarnamError{errorCode: (int)(rc), message: "Varnam Initialization Failed"}
	}
	return &Varnam{handle: v}, nil
}

func (v *Varnam) Learn(text string) *VarnamError {
	rc := C.varnam_learn(v.handle, C.CString(text))
	if rc != 0 {
		errorCode := (int)(rc)
		return &VarnamError{errorCode: errorCode, message: v.getVarnamError(errorCode)}
	}
	return nil
}

func (v *Varnam) getVarnamError(errorCode int) string {
	errormessage := C.varnam_get_last_error(v.handle)
	varnamerrormessage := C.GoString(errormessage)
	return errorMessagesMap[errorCode] + ": " + varnamerrormessage
}
