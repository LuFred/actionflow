package entity

type TimerEntity struct {
	Id                 []byte `db:"id" description:""`
	Name               string `db:"name" description:"名称"`
	Description        string `db:"description" description:"描述"`
	ActionFlowId       []byte `db:"actionFlowId" description:"流程id"`
	ActionFlowVariable string `db:"actionFlowVariable" description:"流程变量"`
	// Cron The scheduling will be based on UTC time
	// ┌───────────── minute (0 - 59)
	// │ ┌───────────── hour (0 - 23)
	// │ │ ┌───────────── day of the month (1 - 31)
	// │ │ │ ┌───────────── month (1 - 12)
	// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
	// │ │ │ │ │
	// │ │ │ │ │
	// * * * * *
	Cron       string `db:"cron" description:"cron表达式"`
	TimerOut   int64  `db:"timerOut" description:"超时时间(s)"`
	Status     string `db:"status" description:"定时器状态"`
	EndTime    int64  `db:"endTime" description:"结束时间"`
	CreatedBy  []byte `db:"createdBy" description:"创建人id"`
	CreatedAt  int64  `db:"createdAt" description:"创建时间"`
	ModifiedBy []byte `db:"modifiedBy" description:"修改人id"`
	ModifiedAt int64  `db:"modifiedAt" description:"修改时间"`
}

func (d *TimerEntity) TableName() string {
	return "timer"
}
