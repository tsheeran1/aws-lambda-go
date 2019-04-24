package chatsess

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

//DATE_FMT tells go the format to use for decoding?
const (
	DATE_FMT = "02-01-2006"
)

//TimetoDB returns a timestamp
func TimetoDB(t time.Time) *string {

	tn := t.Unix()
	return aws.String(strconv.FormatInt(tn, 10))
}

//DBtoTime takes a DB entry and returns a time string
func DBtoTime(s *string) time.Time {
	n, _ := strconv.ParseInt(*s, 10, 64)
	return time.Unix(n, 0)
}
