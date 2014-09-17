package libvarnam

// #cgo pkg-config: varnam
// #include <stdio.h>
// #include <varnam.h>
import "C"

var errorCodes = map[int]string{
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
		return []string{}, &VarnamError{errorCode: (int)(rc), message: "Transliteration Failed: " + errorCodes[(int)(rc)]}
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
		return "", &VarnamError{errorCode: (int)(rc), message: "Reverse Transliteration Failed: " + errorCodes[(int)(rc)]}
	}
	return C.GoString(output), nil
}

func Init(langCode string) (*Varnam, *VarnamError) {
	var v *C.varnam
	var msg *C.char
	rc := C.varnam_init_from_lang(C.CString(langCode), &v, &msg)
	if rc != 0 {
		return nil, &VarnamError{errorCode: (int)(rc), message: C.GoString(msg)}
	}
	return &Varnam{handle: v}, nil
}
}
