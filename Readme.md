# PDF to RAW PRINTING
Convert a pdf file to a kfx, using kpf raw pdf creation with calibre KFX OUTPUT plugin.
The goal is to have an preview in the kindle and other capabilities (drawing, taking notes), locally and on multiple files, even working directly with kindle via MTP.

## Requirements
Golang + dependencies
KFX OUTPUT Plugin

## Program
make the program via make build
$ make build

$ ./build/pdf_raw_printing -pdf /Users/xxx/Downloads/test.pdf -calibre "/Applications/calibre.app/Contents/MacOS/calibre-debug"

$ open test.kfx