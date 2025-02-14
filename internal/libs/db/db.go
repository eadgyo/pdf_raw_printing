package db

import (
	"database/sql"
	"fmt"
	"os"
	"pdf_raw_printing/internal/libs/wion"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db   *sql.DB
	Path string
}

var TESTED = "$ion_symbol_table"

var defaultGCFragment = []string{
	"$ion_symbol_table",
	"book_metadata",
	"book_navigation",
	// "c0",
	// "c0-ad",
	// "c0-spm",
	"content_features",
	// "d6",
	// "d7",
	"document_data",
	// "e9",
	"eidbucket_0",
	"eidbucket_1",
	"eidbucket_10",
	"eidbucket_11",
	"eidbucket_12",
	"eidbucket_13",
	"eidbucket_14",
	"eidbucket_15",
	"eidbucket_16",
	"eidbucket_17",
	"eidbucket_18",
	"eidbucket_19",
	"eidbucket_2",
	"eidbucket_20",
	"eidbucket_21",
	"eidbucket_22",
	"eidbucket_23",
	"eidbucket_24",
	"eidbucket_25",
	"eidbucket_26",
	"eidbucket_27",
	"eidbucket_28",
	"eidbucket_29",
	"eidbucket_3",
	"eidbucket_30",
	"eidbucket_31",
	"eidbucket_32",
	"eidbucket_33",
	"eidbucket_34",
	"eidbucket_35",
	"eidbucket_36",
	"eidbucket_37",
	"eidbucket_38",
	"eidbucket_39",
	"eidbucket_4",
	"eidbucket_40",
	"eidbucket_41",
	"eidbucket_42",
	"eidbucket_43",
	"eidbucket_44",
	"eidbucket_45",
	"eidbucket_46",
	"eidbucket_47",
	"eidbucket_48",
	"eidbucket_49",
	"eidbucket_5",
	"eidbucket_50",
	"eidbucket_51",
	"eidbucket_52",
	"eidbucket_53",
	"eidbucket_54",
	"eidbucket_55",
	"eidbucket_56",
	"eidbucket_57",
	"eidbucket_58",
	"eidbucket_59",
	"eidbucket_6",
	"eidbucket_60",
	"eidbucket_61",
	"eidbucket_62",
	"eidbucket_63",
	"eidbucket_64",
	"eidbucket_65",
	"eidbucket_66",
	"eidbucket_7",
	"eidbucket_8",
	"eidbucket_9",
	// "i4",
	// "i5",
	// "l2",
	"location_map",
	"max_id",
	"metadata",
	"root_entity",
	"rsrc8",
	"yj.kfxid_eid_map",
	"yj.section_pid_count_map",
}

func CreateNewDB(filepath string) (*DB, error) {
	_ = os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	_ = f.Close()

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE capabilities(key char(20), version smallint, primary key (key, version)) without rowid;`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("INSERT INTO capabilities VALUES ($1, $2)", "db.schema", 1)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE fragments(id char(40), payload_type char(10), payload_value blob, primary key (id));`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE gc_reachable(id varchar(40), primary key (id)) without rowid;`)
	if err != nil {
		return nil, err
	}

	for _, el := range defaultGCFragment {
		_, err = db.Exec(`INSERT INTO gc_reachable (id) VALUES ($1);`, el)
		if err != nil {
			return nil, err
		}
	}

	_, err = db.Exec(`CREATE TABLE fragment_properties(id char(40), key char(40), value char(40), primary key (id, key, value)) without rowid;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE gc_fragment_properties(id varchar(40), key varchar(40), value varchar(40), primary key (id, key, value)) without rowid;`)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:   db,
		Path: filepath,
	}, nil
}

func (db *DB) InsertFragment(id string, payloadtype string, payloadvalue []byte) error {
	_, err := db.db.Exec("INSERT INTO fragments (id, payload_type, payload_value) VALUES ($1, $2, $3)", id, payloadtype, payloadvalue)
	return err
}

func (db *DB) InsertHashFragments(id string, payloadType string, v any) error {
	hash24, err := wion.Marshal(v)
	if err != nil {
		return err
	}

	return db.InsertFragment(id, payloadType, hash24)
}

func (db *DB) InsertGCReachable(id string) error {
	_, err := db.db.Exec("INSERT INTO gc_reachable (id) VALUES ($1)", id)
	return err
}

func (db *DB) InsertFragmentProperties(id string, key, value string) error {
	_, err := db.db.Exec(fmt.Sprintf("INSERT INTO fragment_properties (id, key, value) VALUES ('%s', '%s', '%s')", id, key, value))
	return err
}
