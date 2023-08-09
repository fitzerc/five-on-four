package utils

import "strconv"

func UintToString(i uint) string {
    return strconv.FormatUint(uint64(i), 10)
}
