package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"pdf_raw_printing/internal/business"

	"github.com/ledongthuc/pdf"
)

var header = "fa500a5f0100004020ddd75994b7c56c19d8adff59af4dd5a99b8c844d29ea0634cbf4bfe61b9cc4f74d354f9cf73d3aef83c4105892074231f4fcce99dc91e1ff0218575c459b62d523bdfebe3188fcde3047406ec4bc21542fa0c6645559e171ddd75958920742e4dffede15bd362bef83c41034cbf4bf71ddd7591e3b506da70511646e1b5a9df74d354f399e153df3e679a5817be72b94b7c56cf3e679a571ddd759c60baec3b52d2b05f74d354fa705116471ddd759531c714965a85bf5c60baec3707780f26c7b8454c60baec334cbf4bfc60baec3138d854b15bd362b58920742b52d2b051e3b506da84b07ec34cbf4bf0d4d3cd0bb386639d1f64457399e153dd16a935b138d854bf74d354ff74d354fbb386639c60baec3725f3cd3e61b9cc40d4d3cd0a84b07ec94b7c56cf74d354f817be72b71ddd759d1f644573ae34e5ca7051164dd2f7d283ae34e5cf3e679a5ef83c410138d854be25aace0bb386639b52d2b0571ddd759ef83c41071ddd759ed8cc4f7399e153d817be72b531c7149f74d354f71ddd7595892074215bd362ba84b07ec6c7b84546c7b845415bd362b3ed5e703f3e679a578e2cd1df74d354fef83c41034cbf4bf67c37d7c426d947e34cbf4bfb52d2b05e61b9cc499dc91e1426d947e71ddd7596c7b8454ed8cc4f758920742f74d354f65a85bf534cbf4bf531c714934cbf4bf589207427d0f472f725f3cd39cf73d3a67c37d7ca4b451a0725f3cd371ddd75934cbf4bfbb386639f3e679a5abcd7eb5399e153d707780f271ddd759bb38663934cbf4bf3ed5e703138d854b138d854be61b9cc467c37d7ca84b07ec6e1b5a9d426d947e34cbf4bfabcd7eb531f4fcce817be72b707780f2ed8cc4f771ddd759abcd7eb53ed5e703707780f2b52d2b05c60baec3817be72bd1f64457a4b451a0e25aace06e1b5a9dbb3866399cf73d3a58920742e61b9cc4bb3866391e3b506d78e2cd1d7d0f472f531c7149817be72bf3e679a56e1b5a9da70511640d4d3cd015bd362b65a85bf5817be72b138d854bbb386639e61b9cc46c7b845471ddd759abcd7eb565a85bf51e3b506dbb386639abcd7eb5707780f26e1b5a9ddd2f7d2834cbf4bfbb38663978e2cd1d0d4d3cd0c3344a3cc3344a3c71ddd759bb386639abcd7eb57d0f472f6e1b5a9dd16a935b31f4fcce6e1b5a9d99dc91e13ae34e5cbb3866397d0f472fa4b451a0c3344a3cd1f64457707780f2817be72b210e6c8f5892074234cbf4bf399e153ddd2f7d2834cbf4bf399e153da84b07ec3ed5e703e25aace0bb3866393ed5e7036e1b5a9dbb3866396e1b5a9da4b451a015bd362b6c7b8454dd2f7d286e1b5a9d5892074278e2cd1da4b451a0e4dffede9cf73d3aa4b451a031f4fcceb52d2b053ae34e5c3ae34e5c34cbf4bfe4dffede15bd362ba4b45100"

func main() {
	if len(os.Args) == 1 {
		panic("need one pdf in argument")
	}
	wd, _ := os.Getwd()
	cw := wd + "/.tempBook"
	_ = os.RemoveAll(cw)
	err := os.Mkdir(cw, 0777)
	if err != nil {
		panic(err)
	}

	f, r, err := pdf.Open(os.Args[1])
	defer func() { _ = f.Close() }()

	if err != nil {
		panic(err)
	}

	totalPage := r.NumPage()

	err = business.CreateNewPDF(business.PDFInfo{
		Title:         r.Outline().Title,
		Autor:         r.Outline().Title,
		Path:          os.Args[1],
		NumberOfPages: totalPage,
	}, cw)
	if err != nil {
		panic(err)
	}

	headerBytes, _ := hex.DecodeString(header)
	AddMissingSQLiteFile(cw+"/temp.db", headerBytes, cw+"/result.db")
}

func copyDst(src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("Error reading source:", err)
		return
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		fmt.Println("Error writing dest:", err)
		return
	}
}

// Function to add missing sqlite header the SQLite from position 1024 to 2048
func AddMissingSQLiteFile(inputFile string, header []byte, outputFile string) error {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}
	if len(data) < 2048 {
		return fmt.Errorf("file size is less than 2048 bytes")
	}

	modifiedData := make([]byte, len(data)+len(header))
	copy(modifiedData[:1024], data[:1024])
	copy(modifiedData[1024:], header)
	copy(modifiedData[1024+len(header):], data[1024:])

	// Remove bytes from position 1024 to 2048
	return ioutil.WriteFile(outputFile, modifiedData, 0644)
}
