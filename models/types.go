package models

import(
	"github.com/boltdb/bolt"
)

type DBHandler struct {
	DB 				*bolt.DB
	RootBucketName 	[]byte
}

type SuccessResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl string `json:"short_url"`
	Msg string `json:"msg,omitempty"`
}

type FailureResponse struct {
	Error string `json:"error"`
}

type Status struct {
	Error 			error
	SuccessStatus 	*SuccessResponse
	FailureStatus 	*FailureResponse
}

var ALPHABET = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", 
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", 
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", 
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}