package models

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func GetUrl(shortUrl string, DB_Handler *DBHandler) Status {
	var statusResponse Status

	err := DB_Handler.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DB_Handler.RootBucketName)
		innerBucket := bucket.Bucket([]byte(shortUrl))
		if innerBucket == nil {
			fmt.Println("Error: inner bucket returned a nil []byte")
			return fmt.Errorf("No such value for the key \"%s\"", shortUrl)
		}

		//query the original/full URL
		longUrl := innerBucket.Get([]byte("original"))
		if longUrl == nil {
			fmt.Println("bucket.Get() returned nil")
			return fmt.Errorf("innerBucket returned nil...")
		}

		//query the "short" URL
		shortUrl = string(innerBucket.Get([]byte("shortened")))
		//full url queried successfully

		statusResponse = Status{
			Error: nil,
			SuccessStatus: &SuccessResponse{
				OriginalUrl: string(longUrl),
				ShortUrl:    string(shortUrl),
			},
			FailureStatus: nil,
		}
		return nil
	})
	if err != nil {
		statusResponse = Status{
			Error:         err,
			SuccessStatus: nil,
			FailureStatus: &FailureResponse{
					Error: "Invalid URL",
			},
		}
	}
	return statusResponse
}

func SaveUrl(longUrl string, DB_Handler *DBHandler) Status {
	var shortUrl []byte
	var statusResponse Status

	//check that the URL has not been shortened before
	err := DB_Handler.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DB_Handler.RootBucketName)

		// Get the short URL with the given long URL as the key.
		shortUrl = bucket.Get([]byte(longUrl))
		// Get() may return nil, and that means that such key URL
		// has not been saved to the database
		if shortUrl == nil {
			return nil
		} else {
			return fmt.Errorf("The url %s has already been shortened", longUrl)
		}
	})

	// Short the URL only it has not been shortened before.
	// A nil error means that the URL has not been shortened before.
	// So this will happen ONLY if the long URL has not been shortened yet.
	if err == nil {
		// Therefore we proceed to short the url and save it to the database.
		err = DB_Handler.DB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(DB_Handler.RootBucketName)

			// Generate ID for the user.
			// This returns an error only if the Tx is closed or not writeable.
			// That can't happen in an Update() call so I ignore the error check.
			num, _ := bucket.NextSequence()

			// Convert the uint64 to base 62.
			// This function returns a []byte.
			IDByteSlice := base10ToBase62(num)

			// The []byte will be parsed to the string
			// representation of the base 62 number.
			urlID := mapByteSliceToString(IDByteSlice)

			// Create new bucket with the new short url.
			shortUrl = []byte("localhost:8080/api/shorturl/" + urlID)

			innerBucket, err := bucket.CreateBucketIfNotExists([]byte(urlID))
			if err != nil {
				fmt.Println("could not create bucket:", err.Error())
				return err
			}
			fmt.Printf("\nnew bucket for the url ID \"%s\" has been created\n", urlID)
			err = innerBucket.Put([]byte("original"), []byte(longUrl))
			if err != nil {
				fmt.Println("could not insert the url value:", err.Error())
				return err
			}
			err = innerBucket.Put([]byte("shortened"), shortUrl)
			if err != nil {
				fmt.Println("could not insert the shorturl value:", err.Error())
				return err
			}

			// Insert the long URL as the key and the short URL as the value
			// in the root bucket to note that the URL has been saved to
			// the database
			err = bucket.Put([]byte(longUrl), shortUrl)
			if err != nil {
				fmt.Println("could not insert the url value:", err.Error())
				return err
			}
			fmt.Println("database updated successfully")
			return nil
		})
		if err != nil {
			statusResponse = Status{
				Error:			err,
				SuccessStatus: 	nil,
				FailureStatus: 	&FailureResponse{
						Error:	"Internal server error",
				},
			}
		} else {
			statusResponse = Status{
				Error: 			nil,
				SuccessStatus: 	&SuccessResponse{
				OriginalUrl: 	longUrl,
					ShortUrl:	string(shortUrl),
				},
				FailureStatus:	nil,
			}
		}
	} else { //the URL has already been shortened before.
		statusResponse = Status{
			Error: 				nil,
			SuccessStatus: 		&SuccessResponse{
				OriginalUrl:	longUrl,
				ShortUrl:    	string(shortUrl),
				Msg:         	err.Error(),
			},
			FailureStatus:		nil,
		}
	}

	return statusResponse
}

func Start() *DBHandler {
	var DB_Handler = &DBHandler{}
	DB_Handler.RootBucketName = []byte("URL_Index")

	err := DB_Handler.setupDB()
	if err != nil {
		panic(err)
	}
	return DB_Handler
}

func (dbh *DBHandler) setupDB() error {
	db, err := bolt.Open("database/URL_Index.db", 0600, nil)
	if err != nil {
		return fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(dbh.RootBucketName)
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not set up root bucket, %v", err)
	}
	dbh.DB = db
	fmt.Println("DB Setup Done")
	return nil
}
