package business

type AnnotedInteger struct {
	Value      int      `ion:"v"`
	Annotation []string `ion:"annotations"`
}

type AnnotedString struct {
	Value      string   `ion:"v"`
	Annotation []string `ion:"annotations"`
}

type AnnotedFloat struct {
	Value       float64  `ion:"v"`
	Annotations []string `ion:"annotation"`
}

type AnnotedStruct struct {
	Value       interface{} `ion:"v"`
	Annotations []string    `ion:"annotation"`
}

type Annotation int

// --- EIDBUCKET ---
type ContainsElement struct {
	Eid         string `wion:"eid,annotation=kfx_id"`
	SectionName string `wion:"section_name,annotation=kfx_id"`
}

type Eidbucket struct {
	Block      int               `wion:"block"`
	Contains   []ContainsElement `wion:"contains"`
	Annotation Annotation        `wion:"this,annotation=yj.eidhash_eid_section_map"`
}

// --- METADATA ---
type ReadingOrder struct {
	ReadingOrderName Symbol  `wion:"reading_order_name"`
	Sections         []Kfxid `wion:"sections"`
}

type Metadata struct {
	ReadingOrders []ReadingOrder `wion:"reading_orders"`
	Annotation    Annotation     `wion:"this,annotation=metadata"`
}

// --- SECTION ID MAP ---
type ValueMap struct {
	ID         int        `wion:""`
	Reference  string     `wion:",annotation=kfx_id"`
	Annotation Annotation `wion:"this,type=list"`
}

type SectionPositionIdMap struct {
	SectionName string     `wion:"section_name,annotation=kfx_id"`
	Contains    []ValueMap `wion:"contains"`
	Annotation  Annotation `wion:"this,annotation=section_position_id_map"`
}

// --- BOOK_METADATA ---
type BMetadata[T any] struct {
	Key   string `wion:"key"`
	Value T      `wion:"value"`
}

type Ref struct {
	Value      string     `wion:",annotation=kfx_id"`
	Annotation Annotation `wion:"this,type=empty"`
}

type CategorisedMetadata struct {
	Category string `wion:"category"`
	Metadata []any  `wion:"metadata"`
}

type BookMetadata struct {
	CatagoerisedMetadata []CategorisedMetadata `wion:"categorised_metadata"`
	Annotation           Annotation            `wion:"this,annotation=book_metadata"`
}

// --- BookNavigation ---
type NavContainer struct {
	NavType          string     `wion:"nav_type,type=symbol"`
	NavContainerName string     `wion:"nav_container_name,annotation=kfx_id"`
	Entries          []string   `wion:"entries"`
	Annotation       Annotation `wion:"this,annotation=nav_container"`
}

type BookNavigation struct {
	ReadingOrderName string         `wion:"reading_order_name,type=symbol"`
	NavContainers    []NavContainer `wion:"nav_containers"`
}

type BookNavigations struct {
	BookNavigations []BookNavigation `wion:","`
	Annotation      Annotation       `wion:"this,type=empty,annotation=book_navigation"`
}

// --- SECTION cO ---
type Symbol struct {
	Value      string     `wion:",type=symbol"`
	Annotation Annotation `wion:"this,type=empty"`
}

type PageTemplate1 struct {
	Id         string     `wion:"kfx_id,annotation=kfx_id"`
	StoryName  string     `wion:"story_name,annotation=kfx_id"`
	Condition  []Symbol   `wion:"condition,type=sexp"`
	Layout     Symbol     `wion:"layout"`
	TypePage   Symbol     `wion:"type"`
	Annotation Annotation `wion:"this,annotation=structure"`
}

type Width struct {
	Value int    `wion:"value"`
	Unit  Symbol `wion:"unit"`
}

type PageTemplate2 struct {
	Id         string     `wion:"kfx_id,annotation=kfx_id"`
	Width      Width      `wion:"width"`
	StoryName  string     `wion:"story_name,annotation=kfx_id"`
	FixedWidth Width      `wion:"fixed_width"`
	Condition  []Symbol   `wion:"condition,type=sexp"`
	Layout     Symbol     `wion:"layout"`
	TypePage   Symbol     `wion:"type"`
	Annotation Annotation `wion:"this,annotation=structure"`
}

type Section struct {
	SectionName   string     `wion:"section_name,annotation=kfx_id"`
	PageTemplates []any      `wion:"page_templates"`
	Annotation    Annotation `wion:"this,annotation=section"`
}

// --- AuxaliaryData ---
type AuxaliaryData struct {
	Id         string     `wion:"kfx_id,annotation=kfx_id"`
	Metadata   []any      `wion:"metadata"`
	Annotation Annotation `wion:"this,annotation=auxiliary_data"`
}

// --- ModelData ---
type SpecificAuxiliaryData struct {
	Id string `wion:"yj.authoring,annotation=kfx_id"`
}

type Kfxid struct {
	Id         string     `wion:",annotation=kfx_id"`
	Annotation Annotation `wion:"this,type=empty"`
}

type SpecificReadingOrder struct {
	ReadingOrderName Symbol  `wion:"reading_order_name"`
	Sections         []Kfxid `wion:"sections"`
}

type DocumentData struct {
	MaxId         int                    `wion:"max_id"`
	Direction     Symbol                 `wion:"direction"`
	PanZoom       Symbol                 `wion:"pan_zoom"`
	AuxiliaryData SpecificAuxiliaryData  `wion:"auxiliary_data"`
	ReadingOrders []SpecificReadingOrder `wion:"reading_orders"`
	Annotation    Annotation             `wion:"this,annotation=document_data"`
}

// --- ExternalSource ---
type ExternalSource struct {
	MarginLeft     float64    `wion:"margin_left"`
	Format         Symbol     `wion:"format"`
	MarginBottom   float64    `wion:"margin_bottom"`
	MarginRight    float64    `wion:"margin_right"`
	PageIndex      int        `wion:"page_index"`
	Location       string     `wion:"location"`
	AuxiliaryData  Kfxid      `wion:"auxiliary_data"`
	ResourceWidth  float64    `wion:"resource_width"`
	ResourceHeight float64    `wion:"resource_height"`
	ResourceName   Kfxid      `wion:"resource_name"`
	MarginTop      float64    `wion:"margin_top"`
	Annotation     Annotation `wion:"this,annotation=external_resource"`
}

// --- I4 ---
type PageTemplateI4 struct {
	Id          string     `wion:"kfx_id,annotation=kfx_id"`
	FixedWidth  int        `wion:"fixed_width"`
	FixedHeight int        `wion:"fixed_height"`
	FitText     Symbol     `wion:"fit_text"`
	Layout      Symbol     `wion:"layout"`
	Float       Symbol     `wion:"float"`
	Type        Symbol     `wion:"type"`
	ContentList []Kfxid    `wion:"content_list"`
	Annotation  Annotation `wion:"this,annotation=structure"`
}

// --- I5 ---
type PageTemplateI5 struct {
	Id           string     `wion:"kfx_id,annotation=kfx_id"`
	Width        Width      `wion:"width"`
	Height       Width      `wion:"height"`
	Type         Symbol     `wion:"type"`
	ResourceName Kfxid      `wion:"resource_name"`
	Annotation   Annotation `wion:"this,annotation=structure"`
}

// --- StoryLine ---
type StoryLine struct {
	StoryName   Kfxid      `wion:"story_name"`
	ContentList []Kfxid    `wion:"content_list"`
	Annotation  Annotation `wion:"this,annotation=storyline"`
}

// --- MaxID ---
type MaxID struct {
	Value      int        `wion:","`
	Annotation Annotation `wion:"this,type=empty"`
}

// --- yj ---
type YJContains struct {
	SectionName string `wion:"section_name,annotation=kfx_id"`
	Length      int    `wion:"length"`
}

type YJ struct {
	Contains   []YJContains `wion:"contains"`
	Annotation Annotation   `wion:"this,annotation=yj.section_pid_count_map"`
}
