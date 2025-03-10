package account

import (
	"fmt"
	"sync"

	"gitee.com/go-mao/mao/frame"
	"gitee.com/go-mao/mao/libs/try"
)

type Component struct {
	group string
	name  string
	*frame.Taskline
}

func New(t *frame.Taskline, group string, name string) *Component {
	object := new(Component)
	object.Taskline = t
	object.group = group
	object.name = name
	return object
}

func (this *Component) Models() []frame.ModelInterface {
	return []frame.ModelInterface{
		this.NewLog(),
		this.NewAccount(),
	}
}

// 获取，若无则创建
func (this *Component) Get(user interface{}) *AccountModel {
	var accEntity = this.NewAccount()
	var ok bool
	switch data := user.(type) {
	case int64:
		ok = accEntity.Match("UserId", data).Get()
		accEntity.UserId = data
	case UserInterface:
		ok = accEntity.Match("UserId", data.GetUserId()).Get()
		accEntity.UserId = data.GetUserId()
	case *AccountModel:
		ok = data.Id > 0 && data.UserId > 0
		accEntity = data
	default:
		try.Throwf(frame.CODE_WARN, "获取%s用户%s信息失败", this.name, fmt.Sprint(user))
	}
	if !ok && accEntity.UserId > 0 && !accEntity.Create() {
		try.Throwf(frame.CODE_WARN, "创建"+this.group+"失败")
	}
	return accEntity
}

// 获取，若无则创建
func (this *Component) getUserId(user interface{}) int64 {
	switch val := user.(type) {
	case int64:
		return val
	case UserInterface:
		return val.GetUserId()
	case *AccountModel:
		return val.UserId
	default:
		try.Throwf(frame.CODE_WARN, "获取%s用户%s信息失败", this.name, fmt.Sprint(user))
		return 0
	}
}

var accountLock = sync.Map{}

func (this *Component) getLock(userId int64) *sync.Mutex {
	key := fmt.Sprint(this.group, userId)
	val, ok := accountLock.Load(key)
	if !ok {
		val = &sync.Mutex{}
		accountLock.Store(key, val)
	}
	return val.(*sync.Mutex)
}

// 增加余额
func (this *Component) Incr(user interface{}, number int64, typeId int, remarks, note string, handler UserInterface) *AccountModel {
	if number < 0 {
		try.Throwf(frame.CODE_WARN, "操作数量不能小于0")
	}
	if typeId <= 0 {
		try.Throwf(frame.CODE_WARN, "%s余额增加的类型必须大于0", this.name)
	}
	userId := this.getUserId(user)
	lock := this.getLock(userId)
	lock.Lock()
	defer lock.Unlock()
	accEntity := this.Get(user)
	if number == 0 {
		return accEntity
	}
	this.Transaction(func() {
		this.LogDomain().saveLog(accEntity, number, typeId, remarks, note, handler)
		accEntity.Balance += number
		if !accEntity.Cols("Balance").Update() {
			try.Throwf(frame.CODE_SQL, "用户%d%s添加数量操作失败", accEntity.UserId, this.name)
			return
		}
	})
	return accEntity
}

// 减少余额
func (this *Component) Decr(user interface{}, number int64, typeId int, remarks, note string, handler UserInterface) *AccountModel {
	if number < 0 {
		try.Throwf(frame.CODE_WARN, "操作数量不能小于0")
	}
	if typeId >= 0 {
		try.Throwf(frame.CODE_WARN, "%s余额减少的类型必须小于0", this.name)
	}
	userId := this.getUserId(user)
	lock := this.getLock(userId)
	lock.Lock()
	defer lock.Unlock()
	accEntity := this.Get(user)
	if number == 0 {
		return accEntity
	}
	if accEntity.Balance-number < 0 {
		try.Throwf(frame.CODE_WARN, "用户%d%s余额不足", userId, this.name)
		return nil
	}
	this.Transaction(func() {
		this.LogDomain().saveLog(accEntity, 0-number, typeId, remarks, note, handler)
		accEntity.Balance -= number
		if !accEntity.Cols("Balance").Update() {
			try.Throwf(frame.CODE_SQL, "用户%d%s减少数量操作失败", accEntity.UserId, this.name)
			return
		}
	})
	return accEntity
}

// 查看余额是否足够
func (this *Component) CanDecr(user interface{}, number int64) *AccountModel {
	if number < 0 {
		try.Throwf(frame.CODE_WARN, "操作数量不能小于0")
	}
	userId := this.getUserId(user)
	lock := this.getLock(userId)
	lock.Lock()
	defer lock.Unlock()
	accEntity := this.Get(user)
	if accEntity.Balance-number < 0 {
		try.Throwf(frame.CODE_SQL, "用户%d%s余额不足", accEntity.UserId, this.name)
	}
	return accEntity
}

// 统计所有账户余额
func (this *Component) Sum() int64 {
	db := this.OrmSession()
	total, err := db.SumInt(this.NewAccount(), "Balance")
	if err != nil {
		try.Throwf(frame.CODE_SQL, "%s余额统计失败", this.name)
	}
	return total
}

func (this *Component) Find(userId int64, page, limit int, listPtr interface{}) int64 {
	db := this.OrmTable(this.NewAccount())
	if userId > 0 {
		db.Where("`UserId`=?", userId)
	}
	db.Desc("`Id`")
	return this.FindPage(db, listPtr, page, limit)
}

func (this *Component) Delete(userIds []int64) {
	this.Transaction(func() {
		_, err := this.OrmTable(this.NewAccount()).In("UserId", userIds).Delete()
		if err != nil {
			try.Throwf(frame.CODE_SQL, "账户清理失败，ERR:%s", err.Error())
		}
		_, err = this.OrmTable(this.NewLog()).In("UserId", userIds).Delete()
		if err != nil {
			try.Throwf(frame.CODE_SQL, "账户日志清理失败，ERR:%s", err.Error())
		}
	})
}
