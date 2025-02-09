package business

import (
	"pdf_raw_printing/internal/libs/db"
	generator "pdf_raw_printing/internal/libs/idgenerator"
	"strconv"

	"github.com/google/uuid"
)

type PDFInfo struct {
	Title         string
	Autor         string
	Path          string
	NumberOfPages int
}

type PDF struct {
	db         db.DB
	Sections   []string
	Eidbuckets map[int][]KVEid
}

type KVEid struct {
	Key   string
	Value string
}

func NewPdf(tempfolder string) (*PDF, error) {
	myDB, err := db.CreateNewDB(tempfolder + "./temp.db")
	pdf := PDF{
		db:         *myDB,
		Sections:   []string{},
		Eidbuckets: map[int][]KVEid{},
	}

	return &pdf, err
}

func ComputeEID(s string) int {
	value := 0
	for _, c := range s {
		value += int(c)
	}
	return value % 67
}

func (pdf *PDF) AddSectionToEidbucket(key string, value string) {
	eid := ComputeEID(key)
	var k []KVEid
	var exists bool
	if k, exists = pdf.Eidbuckets[eid]; !exists {
		k = []KVEid{}
	}

	k = append(k, KVEid{Key: key, Value: value})
	pdf.Eidbuckets[eid] = k
}

func CreateNewPDF(pdfInfo PDFInfo, tempfolder string) error {
	pdf, err := NewPdf(tempfolder)
	if err != nil {
		return err
	}

	// Start by creating the default infrastructure
	err = pdf.CreateDefaultFragments(pdfInfo)
	if err != nil {
		return err
	}

	// Then create each page
	for i := 0; i < pdfInfo.NumberOfPages; i++ {
		err := pdf.AddPage(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdf *PDF) AddC0(c0 string, c0AD string, l2 string, t1 string, t3 string) error {
	err := pdf.db.InsertFragmentProperties(c0, "child", c0AD)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(c0, "child", l2)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(c0, "child", l2)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(c0, "element_type", "section")
	if err != nil {
		return err
	}

	v := Section{
		SectionName: c0,
		PageTemplates: []any{
			PageTemplate1{
				Id:        t1,
				StoryName: l2,
				Condition: []Symbol{
					{
						Value: "isPortrait",
					},
				},
				Layout: Symbol{
					Value: "vertical",
				},
				TypePage: Symbol{
					Value: "container",
				},
			},
			PageTemplate2{
				Id: t3,
				Width: Width{
					Value: 100,
					Unit: Symbol{
						Value: "percent",
					},
				},
				StoryName: l2,
				FixedWidth: Width{
					Value: 100,
					Unit: Symbol{
						Value: "percent",
					},
				},
				Condition: []Symbol{
					{
						Value: "isLandscape",
					},
				},
				Layout: Symbol{
					Value: "overflow",
				},
				TypePage: Symbol{
					Value: "container",
				},
			},
		},
	}

	return pdf.db.InsertHashFragments(c0, "blob", v)
}

func (pdf *PDF) AddPage(i int) error {
	// c0
	c0 := generator.Generate("c")
	c0AD := c0 + "-ad"
	c0spm := c0 + "-spm"
	l2 := generator.Generate("l")
	d6 := generator.Generate("d")
	d7 := generator.Generate("d")
	e9 := generator.Generate("e")
	i4 := generator.Generate("i")
	i5 := generator.Generate("i")
	_, _ = i4, i5

	t1 := generator.Generate("i")

	t3 := generator.Generate("i")

	pdf.AddC0(c0, c0AD, l2, t1, t3)

	err := pdf.db.InsertFragmentProperties(c0spm, "element_type", "section_position_id_map")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(d6, "element_type", "auxiliary_data")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(d7, "element_type", "auxiliary_data")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("document_data", "child", d7)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(e9, "child", d6)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(e9, "child", "rsrc8")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(e9, "element_type", "external_resource")
	if err != nil {
		return err
	}

	return nil
}

func (pdf *PDF) CreateDefaultFragments(pdfInfo PDFInfo) error {
	err := pdf.db.InsertFragmentProperties("book_metadata", "element_type", "book_metadata")
	if err != nil {
		return err
	}
	err = pdf.AddBookMetadata(pdfInfo.Title, pdfInfo.Autor)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("book_navigation", "element_type", "book_navigation")
	if err != nil {
		return err
	}
	err = pdf.AddBookNavigation()
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("document_data", "element_type", "document_data")
	if err != nil {
		return err
	}
	err = pdf.db.InsertFragmentProperties("max_id", "element_type", "max_id")
	if err != nil {
		return err
	}
	maxId := MaxID{
		Value: 834,
	}
	err = pdf.db.InsertHashFragments("max_id", "blob", maxId)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("metadata", "element_type", "metadata")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("rsrc8", "element_type", "bcRawMedia")
	if err != nil {
		return err
	}
	err = pdf.db.InsertFragment("rsrc8", "path", []byte(pdfInfo.Path))
	if err != nil {
		return err
	}

	for k := range pdf.Eidbuckets {
		err = pdf.db.InsertFragmentProperties("eidbucket_"+strconv.Itoa(k), "element_type", "yj.eidhash_eid_section_map")
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdf *PDF) AddBookNavigation() error {
	bn := BookNavigations{
		BookNavigations: []BookNavigation{
			{
				ReadingOrderName: "default",
				NavContainers: []NavContainer{
					{
						NavType:          "toc",
						NavContainerName: "nA",
						Entries:          []string{},
					},
				},
			},
		},
	}

	return pdf.db.InsertHashFragments("book_navigation", "blob", bn)
}

func (pdf *PDF) AddBookMetadata(title string, autor string) error {
	myuuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	bm := BookMetadata{
		CatagoerisedMetadata: []CategorisedMetadata{
			{
				Category: "kindle_title_metadata",
				Metadata: []any{
					BMetadata[string]{
						Key:   "book_id",
						Value: myuuid.String(),
					},
					BMetadata[string]{
						Key:   "title",
						Value: title,
					},
				},
			},
			{
				Category: "kindle_capability_metadata",
				Metadata: []any{
					BMetadata[int]{
						Key:   "yj_fixed_layout",
						Value: 1,
					},
					BMetadata[int]{
						Key:   "graphical_highlights",
						Value: 1,
					},
					BMetadata[int]{
						Key:   "yj_textbook",
						Value: 1,
					},
				},
			},
			{
				Category: "kindle_ebook_metadata",
				Metadata: []any{
					BMetadata[string]{
						Key:   "selection",
						Value: "enabled",
					},
				},
			},
			{
				Category: "kindle_audit_metadata",
				Metadata: []any{
					BMetadata[string]{
						Key:   "file_creator",
						Value: autor,
					},
					BMetadata[string]{
						Key:   "creator_version",
						Value: "1.93.0.0",
					},
				},
			},
		},
	}
	return pdf.db.InsertHashFragments("book_metadata", "blob", bm)
}
