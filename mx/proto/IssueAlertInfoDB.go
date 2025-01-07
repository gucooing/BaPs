package proto

import (
	"time"
)

type IssueAlertInfoDB struct {
	IssueAlertId   int
	IssueAlertType IssueAlertTypeCode
	StartDate      time.Time
	EndDate        time.Time
	DisplayOrder   []byte
	PublishId      int
	Url            string
	Subject        string
}
