package run

type values struct {
	PUBLIC    string
	PRIVATE   string
	PROTECTED string
	CLASS     string
	NUMBER    string
	VOID      string
	BOOL      string
	REAL      string
	RETURN    string
	MODULE    string
}

// Enum for public use
var Values = &values{
	PUBLIC:    "public: ",
	PRIVATE:   "private: ",
	PROTECTED: "protected: ",
	CLASS:     "class ",
	NUMBER:    "number ",
	BOOL:      "bool ",
	VOID:      "void ",
	REAL:      "real ",
	RETURN:    "return ",
	MODULE:    "module ",
}

var Collections = make(map[string]Collection)

type types struct {
	MAIN     int
	MODULE   int
	VARIABLE int
}

// Enum for public use
var Types = &types{
	MAIN:     1,
	MODULE:   2,
	VARIABLE: 3,
}
