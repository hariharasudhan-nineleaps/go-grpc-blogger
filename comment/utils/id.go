package utils

import (
	"fmt"

	"github.com/teris-io/shortid"
)

func ShortId() string {
	sid, err := shortid.Generate()
	if err != nil {
		panic("Unique Id generation failed")
	}

	fmt.Print("iniq", sid)
	return sid
}
