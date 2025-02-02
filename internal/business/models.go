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
