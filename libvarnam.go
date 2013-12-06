package libvarnam

// #cgo pkg-config: varnam
// #include <stdio.h>
// #include <varnam.h>
import "C"

type Varnam struct {
    handle *C.varnam;
}

type VarnamError struct {
    errorCode int
    message string
}

func (e *VarnamError) Error() string {
    return e.message
}

func (v *Varnam) Transliterate(text string) []string {
    var va *C.varray
    C.varnam_transliterate (v.handle, C.CString (text), &va);
    var i C.int
    var array []string
    for i = 0; i < C.varray_length (va); i++ {
        word := (*C.vword) (C.varray_get (va, i));
        array = append (array, C.GoString (word.text));
    }
    return array
}

func (v *Varnam) ReverseTransliterate(text string) string {
    var output *C.char
    C.varnam_reverse_transliterate (v.handle, C.CString (text), &output);
    return C.GoString (output);
}

func Init(langCode string) (*Varnam, *VarnamError)  {
    var v *C.varnam;
    var msg *C.char
    rc := C.varnam_init_from_lang (C.CString (langCode), &v, &msg)
    if (rc != 0) {
        return nil, &VarnamError {errorCode: (int) (rc), message: C.GoString (msg)}
    }
    return &Varnam {handle: v}, nil
}
