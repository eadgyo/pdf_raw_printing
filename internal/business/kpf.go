package business

import (
	"os"
	"path"

	"github.com/google/uuid"
)

func CreateKCB(tempfolder string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	input := `{
	"book_state" : {
		"book_input_type" : 1,
		"book_manga_comic" : false,
		"book_reading_direction" : 0,
		"book_target_type" : 1,
		"book_virtual_panelmovement" : 0
	},
	"content_hash" : null,
	"metadata" : {
		"book_path" : "resources",
		"edited_tool_versions" : [ "1.93.0.0", "1.96.0.0" ],
		"format" : "yj",
		"global_styling" : true,
		"id" : "` + id.String() + `",
		"platform" : "mac",
		"tool_name" : "KC",
		"tool_version" : "1.96.0.0"
	},
	"tool_data" : {
		"cache_path" : "resources/.cache",
		"created_on" : "2024-Sep-11 18:13:04",
		"last_modified_time" : "2025-Feb-14 22:17:16",
		"link_extract_choice" : false,
		"link_notification_preference" : true
	}
}`

	f, err := os.Create(tempfolder + "/mybook.kcb")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(input)
	if err != nil {
		return err
	}

	return nil
}

func CreateManifestFile(tempfolder string) error {
	manifestfile := `AmazonYJManifest
digital_content_manifest::{
  version:1,
  storage_type:"localSqlLiteDB",
  digital_content_name:"book.kdf"
}`
	f, err := os.Create(tempfolder + "/resources/ManifestFile")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(manifestfile)
	if err != nil {
		return err
	}

	return nil
}

func CreateJournal(tempfolder string) error {
	f, err := os.Create(tempfolder + "/resources/book.kdf-journal")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("")
	if err != nil {
		return err
	}

	return nil
}

func CreateArborescence(pdfpath string, tempfolder1 string) error {
	err := os.Mkdir(path.Join(tempfolder1, "KPF"), 0777)
	if err != nil {
		panic(err)
	}

	tempfolder := path.Join(tempfolder1, "KPF")
	err = CreateKCB(tempfolder)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(tempfolder, "resources"), 0777)
	if err != nil {
		panic(err)
	}
	err = CreateManifestFile(tempfolder)
	if err != nil {
		return err
	}

	err = copyDst(path.Join(tempfolder1, "result.db"), path.Join(tempfolder, "resources", "book.kdf"))
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(tempfolder, "resources", "/res"), 0777)
	if err != nil {
		panic(err)
	}

	err = copyDst(pdfpath, path.Join(tempfolder, "resources", "res", "rsrc8"))
	if err != nil {
		return err
	}
	return nil
}

func copyDst(src string, dst string) error {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
