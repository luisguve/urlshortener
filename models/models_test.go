package models

import (
	"testing"
	"bytes"
	"io/ioutil"
	"os"
	"reflect"

	//"github.com/villegasl/testingutil"
)

type testMappingTable struct {
	Input		[]byte
	Expected	string
}

type testBaseConversionTable struct {
	Input 		uint64
	Expected 	[]byte
}

type DB_Values struct {
	longUrl string
	DB_Handler *DBHandler
}

type SaveOperation struct {
	Input DB_Values
	Expected Status
}

func TestMapByteSliceToString(t *testing.T) {
	tests := []testMappingTable {
		{
			Input:		[]byte{1,2,5},
			Expected:	"abe",
		}, {
			Input:		[]byte{2,1},
			Expected:	"ba",
		}, {
			Input:		[]byte{5,16,47,28,19,57},
			Expected:	"epUBs4",
		}, {
			Input:		[]byte{0},
			Expected:	"",
		}, {
			Input:		[]byte{63},
			Expected:	"ALPHABET overflow",
		}, {
			Input:		[]byte{62},
			Expected:	"9",
		}, 
	}

	var result string

	for _, test := range tests {
		result = mapByteSliceToString(test.Input)
		if result != test.Expected {
			t.Errorf("inputted: %v, expected: %s, received: %s",
				test.Input, test.Expected, result)
		}
	}
}

func TestBase10ToBase62(t *testing.T) {
	tests := []testBaseConversionTable {
		{
			Input:		uint64(125),
			Expected:	[]byte{2,1},
		}, {
			Input:		uint64(754),
			Expected:	[]byte{12,10},
		}, {
			Input:		uint64(10),
			Expected:	[]byte{10},
		}, {
			Input:		uint64(0),
			Expected:	[]byte{0},
		}, {
			Input:		uint64(89),
			Expected:	[]byte{1,27},
		}, {
			Input:		uint64(265748),
			Expected:	[]byte{1,7,8,16},
		}, {
			Input:		uint64(15),
			Expected:	[]byte{15},
		}, {
			Input:		uint64(62),
			Expected:	[]byte{1,0},
		}, {
			Input:		uint64(63),
			Expected:	[]byte{1,1},
		}, 
	}

	var result []byte
	for _, test := range tests {
		result = base10ToBase62(test.Input)
		if !bytes.Equal(result, test.Expected) {
			t.Errorf("inputted: %v, expected: %v, received: %v",
				test.Input, test.Expected, result)
		}
	}
}

func TestSaveURL(t *testing.T) {
    db := NewTestDB()
    defer db.Close()

    tests := []SaveOperation {
    	{
    		Input: DB_Values{
	    		longUrl: 	"www.freecodecamp.org",
	    		DB_Handler:	&DBHandler{db.DB, db.RootBucketName},
    		},
    		Expected: Status{
    			Error: 			nil,
				SuccessStatus: 	&SuccessResponse{
				OriginalUrl: 	"www.freecodecamp.org",
				ShortUrl:		"localhost:8080/api/shorturl/a",
				},
				FailureStatus:	nil,
    		},
    	}, {
    		Input: DB_Values{
	    		longUrl: 	"www.wikipedia.org",
	    		DB_Handler:	&DBHandler{db.DB, db.RootBucketName},
    		},
    		Expected: Status{
    			Error: 			nil,
				SuccessStatus: 	&SuccessResponse{
				OriginalUrl: 	"www.wikipedia.org",
				ShortUrl:		"localhost:8080/api/shorturl/b",
				},
				FailureStatus:	nil,
    		},
    	}, {
    		Input: DB_Values{
	    		longUrl: 	"www.freecodecamp.org",
	    		DB_Handler:	&DBHandler{db.DB, db.RootBucketName},
    		},
    		Expected: Status{
    			Error: 			nil,
				SuccessStatus: 	&SuccessResponse{
				OriginalUrl: 	"www.freecodecamp.org",
				ShortUrl:		"localhost:8080/api/shorturl/a",
				Msg:			"The url www.freecodecamp.org has already been shortened",
				},
				FailureStatus:	nil,
    		},
    	},
    }

    var status Status
    for _, test := range tests {
    	status = SaveUrl(test.Input.longUrl, test.Input.DB_Handler)
    	if !reflect.DeepEqual(status, test.Expected) {
    		t.Errorf("inputted: %v, expected: %v, received: %v",
    			test.Input, test.Expected, status)
    	}
    }
}

// NewTestDB returns a TestDB using a temporary path.
func NewTestDB() *TestDB {
    // Retrieve a temporary path.
    f, err := ioutil.TempFile("", "")
    if err != nil {
        panic("temp file: " + err.Error())
    }
    path := f.Name()
    f.Close()
    os.Remove(path)
    // Open the database.
    db, err := Open(path, []byte("URL_Index"), 0600)
    if err != nil {
		panic("open: " + err.Error())
    }
    // Return wrapped type.
    return &TestDB{db}
}

// Close and delete Bolt database.
func (db *TestDB) Close() {
    defer os.Remove(db.DB.Path())
    db.DB.Close()
}