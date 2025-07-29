package version

import (
	"fmt"
	"time"
)

var Version string = "dev"
var BuildTime string

func init() {
	if BuildTime == "" {
		BuildTime = fmt.Sprint(time.Now().Unix())
	}
}
