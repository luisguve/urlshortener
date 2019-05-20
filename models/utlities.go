package models

func mapUint64SliceToString(digits []uint64) string {
	var urlID string

    for _, v := range digits {
    // Each value within the digits slice must be decremented
    // in order to get the right letter of the ALPHABET array
    	v--
     	urlID += ALPHABET[v]
    }
    return urlID
}

func base10ToBase62(num uint64) []uint64 {
	var digits []uint64
    var remainder uint64

    // Populate the digits slice with the corresponding 
    // indexes to map the ID using the ALPHABET array
    for num > 0 {
    	remainder = num % uint64(62)
       	digits = append(digits, remainder)
       	num = num / uint64(62)
    }

    // The slice must be reversed in order to have be valid base 62 number
    reverseSlice(digits)

    return digits
}

func reverseSlice(digits []uint64) {
	for i := 0; i < len(digits) / 2; i++ {
		j := len(digits) - i - 1
		digits[i], digits[j] = digits[j], digits[i]
	}
}