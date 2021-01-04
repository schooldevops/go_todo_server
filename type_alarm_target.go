package main

type alarmTarget struct {
	ID          int64  `json:"id,omitempty"`
	TodoAlarmID int64  `json:"todo_alarm_id,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	ActiveYn    string `json:"active_yn,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	ModifiedAt  string `json:"modified_at,omitempty"`
}
