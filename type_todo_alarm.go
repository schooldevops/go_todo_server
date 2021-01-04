package main

type todoAlarm struct {
	ID         int64  `json:"id,omitempty"`
	TodoID     int64  `json:"todo_id,omitempty"`
	PeriodType string `json:"period_type,omitempty"`
	AlarmType  string `json:"alarm_type,omitempty"`
	AlarmDate  string `json:"alarm_date,omitempty"`
	AlarmTime  string `json:"alarm_time,omitempty"`
	ActiveYn   string `json:"active_yn,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	ModifiedAt string `json:"modified_at,omitempty"`
}
