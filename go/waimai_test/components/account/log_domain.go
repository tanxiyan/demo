package account

import (
	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

type logDomain struct {
	*Component
}

func (this *Component) LogDomain() *logDomain {
	return &logDomain{this}
}

// 金额统计
func (this *logDomain) Sum(userId int64, classIds ...interface{}) int64 {
	db := this.OrmSession()
	if userId > 0 {
		db.Where("`UserId`=?", userId)
	}
	if len(classIds) > 0 {
		db.In("TypeId", classIds)
	}
	total, err := db.SumInt(this.NewLog(), "Number")
	if err != nil {
		try.Throwf(frame.CODE_SQL, "用户%d%s日志统计失败：%s", userId, this.name, err.Error())
	}
	return total
}

// 保存流水
func (this *logDomain) saveLog(account *AccountModel, number int64, typeId int, remarks, note string, handler UserInterface) {
	logEntity := this.NewLog()
	logEntity.TypeId = typeId
	logEntity.UserId = account.UserId
	logEntity.BalanceBefore = account.Balance
	logEntity.BalanceAfter = account.Balance + number
	logEntity.Remarks = remarks
	logEntity.Note = note
	logEntity.Number = number
	if handler != nil {
		logEntity.CreateUserId = handler.GetUserId()
		logEntity.CreateUserGroup = handler.GetUserGroup()
	}
	if !logEntity.Create() {
		try.Throwf(frame.CODE_SQL, "用户%d%s日志保存失败", account.Id, this.name)
		return
	}
}

type FindLogArgs struct {
	UserId  int64
	TypeId  int
	BeginAt string
	EndAt   string
}

// 流水查询
func (this *logDomain) Find(args FindLogArgs, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewLog())
	if args.UserId > 0 {
		db.Where("`UserId`=?", args.UserId)
	}
	if args.TypeId != 0 {
		db.Where("`TypeId`=?", args.TypeId)
	}
	if args.BeginAt != "" {
		db.And("`CreateTime`>?", args.BeginAt)
	}
	if args.EndAt != "" {
		db.And("`CreateTime`<?", args.EndAt)
	}
	db.Desc("`Id`")
	return this.FindPage(db, listPtr, page, limit)
}

func (this *logDomain) Delete(beginAt, endAt string, ids []int64) {
	logsInfo := this.NewLog()
	db := this.OrmTable(logsInfo)
	if beginAt != "" && endAt != "" {
		db.And("`CreateTime` between ? and ?", beginAt, endAt)
	}
	if len(ids) > 0 {
		db.In("Id", ids)
	}
	_, err := db.Delete()
	if err != nil {
		try.Throwf(frame.CODE_SQL, "%s日志删除失败：%s", this.name, err.Error())
	}
}
