package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeOperation(t *testing.T) {
	timeZone := 8
	loc := time.FixedZone("ZONE", timeZone*60*60)
	nowTime := time.Now().Unix()
	resetSec := int64(15 * 60 * 60)
	today := time.Unix(nowTime-resetSec, 0).In(loc).Format("2006-01-02")
	todayTime, _ := time.ParseInLocation("2006-01-02", today, loc)
	nextDayTime := todayTime.AddDate(0, 0, 1)
	nextResetTime := nextDayTime.Add(time.Second * time.Duration(resetSec))
	fmt.Println(nextResetTime)
}
