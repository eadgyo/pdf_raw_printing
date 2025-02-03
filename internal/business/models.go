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
	ReadingOrderName string   `wion:"reading_order_name"`
	Sections         []string `wion:"sections"`
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
