package business

import (
	"encoding/hex"
	"pdf_raw_printing/internal/libs/db"
	generator "pdf_raw_printing/internal/libs/idgenerator"
	"strconv"

	"github.com/google/uuid"
)

var DEBUG_ONE_PAGE = false

var ion_symbol_table = "e00100eaeea08183de9c8822034286be95de93848a594a5f73796d626f6c7385210a88220339"

type PDFInfo struct {
	Title         string
	Autor         string
	Path          string
	Id            string
	NumberOfPages int
}

type PDF struct {
	db         db.DB
	Sections   []string
	Eidbuckets map[int][]KVEid
	d6         string
	d7         string
}

type KVEid struct {
	Key   string
	Value string
}

func NewPdf(tempfolder string) (*PDF, error) {
	myDB, err := db.CreateNewDB(tempfolder + "/temp.db")

	if err != nil {
		return nil, err
	}

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

	// Start by creating the init
	generator.Register("d6")
	generator.Register("d7")
	pdf.d6 = "d6"
	pdf.d7 = "d7"

	// Then create each page
	for i := 0; i < pdfInfo.NumberOfPages; i++ {
		err := pdf.AddPage(i)
		if err != nil {
			return err
		}

		if DEBUG_ONE_PAGE {
			if i == 0 {
				break
			}
		}
	}

	err = pdf.CreateDefaultFragments(pdfInfo)
	if err != nil {
		return err
	}

	return nil
}

func (pdf *PDF) AddD6(d6 string, path string) error {
	err := pdf.db.InsertFragmentProperties(d6, "element_type", "auxiliary_data")
	if err != nil {
		return err
	}

	v := AuxaliaryData{
		Id: d6,
		Metadata: []any{
			BMetadata[string]{
				Key:   "type",
				Value: "resource",
			},
			BMetadata[string]{
				Key:   "resource_stream",
				Value: "rsrc8",
			},
			BMetadata[string]{
				Key:   "size",
				Value: "12844",
			},
			BMetadata[string]{
				Key:   "modified_time",
				Value: "1725439015",
			},
			BMetadata[string]{
				Key:   "location",
				Value: path,
			},
		},
	}

	return pdf.db.InsertHashFragments(d6, "blob", v)
}

func (pdf *PDF) AddE9(e9 string, pageIndex int) error {
	err := pdf.db.InsertFragmentProperties(e9, "child", pdf.d6)
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

	v := ExternalSource{
		MarginLeft: 0,
		Format: Symbol{
			Value: "pdf",
		},
		MarginBottom: 0,
		MarginRight:  0.5,
		PageIndex:    pageIndex,
		Location:     "rsrc8",
		AuxiliaryData: Kfxid{
			Id: pdf.d6,
		},
		ResourceWidth:  596,
		ResourceHeight: 842,
		ResourceName: Kfxid{
			Id: e9,
		},
		MarginTop: 0,
	}

	return pdf.db.InsertHashFragments(e9, "blob", v)
}

func (pdf *PDF) AddD7(d7 string) error {
	err := pdf.db.InsertFragmentProperties(d7, "element_type", "auxiliary_data")
	if err != nil {
		return err
	}

	v := AuxaliaryData{
		Id: d7,
		Metadata: []any{
			BMetadata[[]Ref]{
				Key: "auxData_resource_list",
				Value: []Ref{
					{
						Value: pdf.d6,
					},
				},
			},
		},
	}

	return pdf.db.InsertHashFragments(d7, "blob", v)
}

func (pdf *PDF) AddC0Spm(c0 string, c0spm string, t1 string, t3 string, i4 string, i5 string) error {
	err := pdf.db.InsertFragmentProperties(c0spm, "element_type", "section_position_id_map")
	if err != nil {
		return err
	}

	pdf.AddSectionToEidbucket(t1, c0)
	pdf.AddSectionToEidbucket(t3, c0)

	v := SectionPositionIdMap{
		Contains: []ValueMap{
			{
				ID:        1,
				Reference: t1,
			},
			{
				ID:        2,
				Reference: i4,
			},
			{
				ID:        3,
				Reference: i5,
			},
			{
				ID:        4,
				Reference: t3,
			},
		},
		SectionName: c0,
	}

	return pdf.db.InsertHashFragments(c0spm, "blob", v)
}

func (pdf *PDF) AddC0(c0 string, c0AD string, l2 string, t1 string, t3 string) error {
	pdf.AddSectionToEidbucket(c0, c0)
	pdf.Sections = append(pdf.Sections, c0)

	err := pdf.db.InsertFragmentProperties(c0, "child", c0AD)
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

func (pdf *PDF) AddI5(c0 string, i5 string, e9 string) error {
	err := pdf.db.InsertFragmentProperties(i5, "child", e9)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(i5, "element_type", "structure")
	if err != nil {
		return err
	}

	pdf.AddSectionToEidbucket(i5, c0)

	v := PageTemplateI5{
		Id: i5,
		Width: Width{
			Value: 100,
			Unit: Symbol{
				Value: "percent",
			},
		},
		Height: Width{
			Value: 100,
			Unit: Symbol{
				Value: "percent",
			},
		},
		Type: Symbol{
			Value: "image",
		},
		ResourceName: Kfxid{
			Id: e9,
		},
	}

	return pdf.db.InsertHashFragments(i5, "blob", v)

}

func (pdf *PDF) AddI4(c0 string, i4 string, i5 string) error {
	err := pdf.db.InsertFragmentProperties(i4, "child", i5)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(i4, "element_type", "structure")
	if err != nil {
		return err
	}

	pdf.AddSectionToEidbucket(i4, c0)

	v := PageTemplateI4{
		Id:          i4,
		FixedWidth:  59600,
		FixedHeight: 84200,
		FitText: Symbol{
			Value: "force",
		},
		Layout: Symbol{
			Value: "scale_fit",
		},
		Float: Symbol{
			Value: "center",
		},
		Type: Symbol{
			Value: "container",
		},
		ContentList: []Kfxid{
			{
				Id: i5,
			},
		},
	}

	return pdf.db.InsertHashFragments(i4, "blob", v)
}

func (pdf *PDF) AddL2(l2 string, i4 string) error {
	err := pdf.db.InsertFragmentProperties(l2, "child", l2)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(l2, "child", i4)
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties(l2, "element_type", "storyline")
	if err != nil {
		return err
	}

	v := StoryLine{
		StoryName: Kfxid{
			Id: l2,
		},
		ContentList: []Kfxid{
			{
				Id: i4,
			},
		},
	}

	return pdf.db.InsertHashFragments(l2, "blob", v)
}

func (pdf *PDF) AddC0AD(c0AD string) error {
	err := pdf.db.InsertFragmentProperties(c0AD, "element_type", "auxiliary_data")
	if err != nil {
		return err
	}

	v := AuxaliaryData{
		Id: c0AD,
		Metadata: []any{
			BMetadata[int]{
				Key:   "page_rotation",
				Value: 0,
			},
		},
	}

	return pdf.db.InsertHashFragments(c0AD, "blob", v)
}

func (pdf *PDF) AddPage(i int) error {
	// c0
	c0 := generator.Generate("c")
	c0AD := c0 + "-ad"
	c0spm := c0 + "-spm"
	generator.Register(c0spm)
	generator.Register(c0AD)
	l2 := generator.Generate("l")
	e9 := generator.Generate("e")
	i4 := generator.Generate("i")
	i5 := generator.Generate("i")

	t1 := generator.Generate("t")

	t3 := generator.Generate("t")

	if DEBUG_ONE_PAGE {
		t1 = "t1"
		i5 = "i5"
		i4 = "i4"
		l2 = "l2"
		c0 = "c0"
		e9 = "e9"
		t3 = "t3"
		c0AD = c0 + "-ad"
		c0spm = c0 + "-spm"
	}

	err := pdf.AddC0(c0, c0AD, l2, t1, t3)
	if err != nil {
		return err
	}

	err = pdf.AddC0AD(c0AD)
	if err != nil {
		return err
	}

	err = pdf.AddC0Spm(c0, c0spm, t1, t3, i4, i5)
	if err != nil {
		return err
	}

	err = pdf.AddE9(e9, i)
	if err != nil {
		return err
	}

	err = pdf.AddI5(c0, i5, e9)
	if err != nil {
		return err
	}

	err = pdf.AddI4(c0, i4, i5)
	if err != nil {
		return err
	}

	err = pdf.AddL2(l2, i4)
	if err != nil {
		return err
	}

	return nil
}

func (pdf *PDF) AddMetadata() error {
	err := pdf.db.InsertFragmentProperties("metadata", "element_type", "metadata")
	if err != nil {
		return err
	}

	sections := []Kfxid{}

	for _, v := range pdf.Sections {
		sections = append(sections, Kfxid{Id: v})
	}

	v := Metadata{
		ReadingOrders: []ReadingOrder{
			{
				ReadingOrderName: Symbol{
					Value: "default",
				},
				Sections: sections,
			},
		},
	}

	return pdf.db.InsertHashFragments("metadata", "blob", v)
}

func (pdf *PDF) AddMaxId() error {
	err := pdf.db.InsertFragmentProperties("max_id", "element_type", "max_id")
	if err != nil {
		return err
	}
	maxId := MaxID{
		Value: 834,
	}
	return pdf.db.InsertHashFragments("max_id", "blob", maxId)
}

func (pdf *PDF) AddEidBuckets() error {
	for id, els := range pdf.Eidbuckets {
		contains := []ContainsElement{}
		for _, v := range els {
			contains = append(contains, ContainsElement{Eid: v.Key, SectionName: v.Value})
		}

		v := Eidbucket{
			Block:    id,
			Contains: contains,
		}

		err := pdf.db.InsertHashFragments("eidbucket_"+strconv.Itoa(id), "blob", v)
		if err != nil {
			return err
		}

		err = pdf.db.InsertFragmentProperties("eidbucket_"+strconv.Itoa(id), "element_type", "yj.eidhash_eid_section_map")
		if err != nil {
			return err
		}
	}
	return nil
}

func (pdf *PDF) AddRoot() error {
	err := pdf.db.InsertFragmentProperties("$ion_symbol_table", "element_type", "$ion_symbol_table")
	if err != nil {
		return err
	}
	bytesIon, _ := hex.DecodeString(ion_symbol_table)
	return pdf.db.InsertFragment("$ion_symbol_table", "blob", bytesIon)
}

func (pdf *PDF) CreateDefaultFragments(pdfInfo PDFInfo) error {
	err := pdf.AddRoot()
	if err != nil {
		return err
	}

	err = pdf.AddBookMetadata(pdfInfo.Title, pdfInfo.Autor)
	if err != nil {
		return err
	}

	err = pdf.AddBookNavigation()
	if err != nil {
		return err
	}

	err = pdf.AddDocumentData()
	if err != nil {
		return err
	}

	err = pdf.AddD6(pdf.d6, pdfInfo.Path)
	if err != nil {
		return err
	}

	err = pdf.AddD7(pdf.d7)
	if err != nil {
		return err
	}

	err = pdf.AddMaxId()
	if err != nil {
		return err
	}

	err = pdf.AddMetadata()
	if err != nil {
		return err
	}

	err = pdf.AddSectionPidCountMap()
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

	err = pdf.AddEidBuckets()
	if err != nil {
		return err
	}

	return nil
}

func (pdf *PDF) AddSectionPidCountMap() error {

	err := pdf.db.InsertFragmentProperties("yj.section_pid_count_map", "element_type", "yj.section_pid_count_map")
	if err != nil {
		return err
	}
	contains := []YJContains{}

	for _, v := range pdf.Sections {
		contains = append(contains, YJContains{
			SectionName: v,
			Length:      4,
		})
	}

	v := YJ{
		Contains: contains,
	}

	return pdf.db.InsertHashFragments("yj.section_pid_count_map", "blob", v)
}

func (pdf *PDF) AddDocumentData() error {
	err := pdf.db.InsertFragmentProperties("document_data", "element_type", "document_data")
	if err != nil {
		return err
	}

	err = pdf.db.InsertFragmentProperties("document_data", "child", pdf.d7)
	if err != nil {
		return err
	}

	sections := []Kfxid{}
	for _, v := range pdf.Sections {
		sections = append(sections, Kfxid{Id: v})
	}

	v := DocumentData{
		MaxId: generator.GetSize(),
		Direction: Symbol{
			Value: "ltr",
		},
		PanZoom: Symbol{
			Value: "enabled",
		},
		AuxiliaryData: SpecificAuxiliaryData{
			Id: pdf.d7,
		},
		ReadingOrders: []SpecificReadingOrder{
			{
				ReadingOrderName: Symbol{
					Value: "default",
				},
				Sections: sections,
			},
		},
	}

	return pdf.db.InsertHashFragments("document_data", "blob", v)
}

func (pdf *PDF) AddBookNavigation() error {
	err := pdf.db.InsertFragmentProperties("book_navigation", "element_type", "book_navigation")
	if err != nil {
		return err
	}

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
	err := pdf.db.InsertFragmentProperties("book_metadata", "element_type", "book_metadata")
	if err != nil {
		return err
	}

	myuuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	myuuidstr := myuuid.String()

	if DEBUG_ONE_PAGE {
		myuuidstr = "dNdrKhVjR26t_cZ_uHYkOA0"
		title = "ttt"
		autor = "KC"
	}

	bm := BookMetadata{
		CatagoerisedMetadata: []CategorisedMetadata{
			{
				Category: "kindle_title_metadata",
				Metadata: []any{
					BMetadata[string]{
						Key:   "book_id",
						Value: myuuidstr,
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
