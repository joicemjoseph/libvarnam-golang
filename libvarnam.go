package libvarnam

// #cgo pkg-config: varnam
// #include <stdio.h>
// #include <varnam.h>
import "C"
import "fmt"

func Init() {
    var v *C.varnam;
    var msg *C.char
    var va *C.varray
    var word *C.vword

    rc := C.varnam_init_from_lang (C.CString ("ml"), &v, &msg)
    rc = C.varnam_transliterate (v, C.CString ("navaneeth"), &va);
    var i C.int
    for i = 0; i < C.varray_length (va); i++ {
        word = (*C.vword) (C.varray_get (va, i));
        fmt.Println (C.GoString (word.text));
    }
    fmt.Println (rc);
}
