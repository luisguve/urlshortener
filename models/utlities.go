package models

// mapUint64SliceToString takes in a byte slice and returns
// a string containing the alphabetical representation, i.e [a-z-A-Z-0-9],
// of each base 62-encoded digit in the given slice.
// 
// Given the byte slice {1,2,5} it should return the string "abe"
func mapByteSliceToString(digits []byte) string {
	var urlID string

    for _, v := range digits {
        // 62 is the limit, but len(ALPHABET) == 63.
        if int(v) > len(ALPHABET) - 1 {
            return "ALPHABET overflow"
        }
     	urlID += ALPHABET[v]
    }
    return urlID
}

//  base10ToBase62 takes in a uint64 number and returns a slice
//  of uint64 containing each digit of the uint64 number base 62-encoded.
//
//  Given the uint64 number 125, it should return the uint64 slice [2,1]
func base10ToBase62(num uint64) []byte {
	var digitsResult []byte
    var remainder uint64

    if num == 0 {
        return []byte{0}
    }

    // Populate the digitsResult slice with the corresponding 
    // indexes to map the ID using the ALPHABET array
    for num > 0 {
    	remainder = num % uint64(62)
       	digitsResult = append(digitsResult, byte(remainder))
       	num = num / uint64(62)
    }

    // The slice must be reversed in order to have be valid base 62 number
    reverseSlice(digitsResult)

    return digitsResult
}

func reverseSlice(digits []byte) {
	for i := 0; i < len(digits) / 2; i++ {
		j := len(digits) - i - 1
		digits[i], digits[j] = digits[j], digits[i]
	}
}