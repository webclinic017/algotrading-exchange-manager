package kite

import (
	"strconv"
)

func GetSymbols() {

	// call api to load the instruments CSV file into db

	// get the instruments list from db

	var i = 0
	var val uint64
	for _, value := range InsNamesMap {
		val, _ = strconv.ParseUint(value, 10, 64)
		Tokens = append(Tokens, uint32(val))
		i++
	}

}
