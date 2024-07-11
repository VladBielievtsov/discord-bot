package utils

import "log"

func ErrorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Colors = []int{
	0x19a5ff,
	0xfee719,
	0xfe4f19,
	0xf00c18,
	0xff1791,
	0x04e762,
	0x8900f2,
}
