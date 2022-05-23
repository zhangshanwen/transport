package model

import (
	"errors"
	"time"
)

type (
	Status int8
	Task   struct {
		BaseModel
		Name       string `json:"name"`        // 任务名称
		Remark     string `json:"remark"`      // 备注
		Status     Status `json:"status"`      // 状态  -2:过期 -1:停止 0:待生效 1:生效中
		Cmd        string `json:"cmd"`         // 任务指令
		EffectDate int64  `json:"effect_date"` // 生效时间
		ExpiryDate int64  `json:"expiry_date"` // 失效时间
		StartAt    int    `json:"start_at"`    // 任务开始执行时间
		EndAt      int    `json:"end_at"`      // 任务结束时间
	}
	TaskDevice struct {
		TaskId   int64 `json:"task_id"`
		DeviceId int64 `json:"device_id"`
	}
)

const (
	StatusExpiry Status = iota - 1
	StatusStop
	StatusIdle
	StatusRunning
)

func (t *Task) Verify() (err error) {
	t.Status = StatusIdle

	now := time.Now().Unix()
	if t.EffectDate >= t.ExpiryDate {
		return errors.New("过期时间不得小于等于生效时间")
	}
	if t.ExpiryDate-t.EffectDate < 60*60*24*1 {
		// 时间间隔
		return errors.New("时间间隔不得小于1天")
	}
	if t.ExpiryDate+60*60*24*1 <= now {
		return errors.New("过期时间不得小于1天")
	}
	if now >= t.ExpiryDate {
		return errors.New("过期时间不得小于当前时间")
	}
	if now >= t.EffectDate {
		t.Status = StatusRunning
	}
	return
}
