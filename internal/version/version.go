package version

import (
	"fmt"
	"log"
	"time"
)

var Version string = "dev"
var BuildTime string

func init() {
	log.Print(Version)
	log.Print(BuildTime)
	if BuildTime == "" {
		BuildTime = fmt.Sprint(time.Now().Unix())
	}
}
