# PDF to RAW PRINTING
Convert a pdf file to a kindle kfx, using kpf raw pdf creation with calibre KFX OUTPUT plugin.
The goal is to have a preview in the kindle and other capabilities (drawing, taking notes), locally and on multiple files, even working directly with kindle via MTP.

## Requirements
Golang + dependencies
KFX OUTPUT Plugin

## Program
make the program via make build
```
$ make build
```

### One file
```
$ ./build/pdf_raw_printing -pdf /Users/xxx/Downloads/test.pdf -calibre "/Applications/calibre.app/Contents/MacOS/calibre-debug"

$ open test.kfx
```


### Folder
```
$ ./build/pdf_raw_printing -pdf /Users/xxx/Downloads/tests -dest /Users/xxx/Downloads/tests -calibre "/Applications/calibre.app/Contents/MacOS/calibre-debug"

$ ls /Users/xxx/Downloads/tests
file1.pdf
file1.kfx
file2.pdf
file2.kfx
```

