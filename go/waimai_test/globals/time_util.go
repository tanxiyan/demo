package globals

const (
	TimeStart = " 00:00:00"
	TimeEnd   = " 23:59:59"
)

func FormatStartTime(t string) string {
	//判断是否只有yyyy-mm-dd的日期格式
	if len(t) == 10 {
		return t + TimeStart
	} else {
		return t
	}
}
func FormatEndTime(t string) string {
	//判断是否只有yyyy-mm-dd的日期格式
	if len(t) == 10 {
		return t + TimeEnd
	} else {
		return t
	}
}
