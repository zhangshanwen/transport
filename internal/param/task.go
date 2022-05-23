package param

type (
	TaskCreate struct {
		Name       string `json:"name"`
		EffectDate int64  `json:"effect_date"` // 生效时间
		ExpiryDate int64  `json:"expiry_date"` // 失效时间
		StartAt    int    `json:"start_at"`    // 任务开始执行时间
		EndAt      int    `json:"end_at"`      // 任务结束时间
	}
	Task struct {
		Pagination
	}
)
