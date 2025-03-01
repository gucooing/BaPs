package proto

type IssueAlertTypeCode int32

const (
	IssueAlertTypeCode_All                  = 1
	IssueAlertTypeCode_File_Target          = 2
	IssueAlertTypeCode_AllButFile_Exception = 3
)

var (
	IssueAlertTypeCode_name = map[int32]string{
		1: "All",
		2: "File_Target",
		3: "AllButFile_Exception",
	}
	IssueAlertTypeCode_value = map[string]int32{
		"All":                  1,
		"File_Target":          2,
		"AllButFile_Exception": 3,
	}
)

func (x IssueAlertTypeCode) String() string {
	return IssueAlertTypeCode_name[int32(x)]
}
