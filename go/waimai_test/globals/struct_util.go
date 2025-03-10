package globals

import (
	"reflect"
)

// CopyStruct 将结构体source复制给dst，只复制相同名称和相同类型的
// CopyStruct(a,b)  a可以传值，引用，b只能引用
func CopyStruct(src, dst interface{}, cols ...string) interface{} {
	st := reflect.TypeOf(src)
	sv := reflect.ValueOf(src)
	dt := reflect.TypeOf(dst)
	dv := reflect.ValueOf(dst)
	if st.Kind() == reflect.Ptr { //处理指针
		st = st.Elem()
		sv = sv.Elem()
	}
	if dt.Kind() == reflect.Ptr { //处理指针
		dt = dt.Elem()
	}
	if st.Kind() != reflect.Struct || dt.Kind() != reflect.Struct { //如果不是struct类型，直接返回dst
		return dst
	}
	dv = reflect.ValueOf(dv.Interface())
	// 遍历TypeOf 类型
	for i := 0; i < dt.NumField(); i++ { //通过索引来取得它的所有字段，这里通过t.NumField来获取它多拥有的字段数量，同时来决定循环的次数
		f := dt.Field(i) //通过这个i作为它的索引，从0开始来取得它的字段
		dVal := dv.Elem().Field(i)
		sVal := sv.FieldByName(f.Name)
		//fmt.Println(dVal.CanSet())
		//src数据有效，且dst字段能赋值,类型一致
		if sVal.IsValid() && dVal.CanSet() && f.Type.Kind() == sVal.Type().Kind() {
			if len(cols) > 0 {
				if InArrayStr(f.Name, cols...) {
					if sVal.Type().Kind() == reflect.Slice {
						continue
					}
					dVal.Set(sVal)
				}
			} else {
				if sVal.Type().Kind() == reflect.Slice {
					continue
				}
				dVal.Set(sVal)
			}
		}
	}
	return dst
}

func GenSlice(req interface{}) interface{} {

	typ := reflect.TypeOf(req)

	// 通过类型 产生切片类型
	sliceType := reflect.SliceOf(typ)

	// 创建切片
	sliceVal := reflect.MakeSlice(sliceType, 0, 0)

	// 赋值
	vals := reflect.Append(sliceVal, reflect.ValueOf(req))

	return vals.Interface()
}
func InArrayStr(name string, cols ...string) bool {
	isHave := false
	for _, col := range cols {
		if name == col {
			isHave = true
			break
		}
	}
	return isHave
}

func InArray(name int64, cols ...int64) bool {
	isHave := false
	for _, col := range cols {
		if name == col {
			isHave = true
			break
		}
	}
	return isHave
}

func GetFieldList(bean interface{}) []string {
	fields := make([]string, 0)
	beanType := reflect.TypeOf(bean)
	beanValue := reflect.ValueOf(bean)
	if beanType.Kind() == reflect.Ptr { //处理指针
		beanType = beanType.Elem()
		beanValue = beanValue.Elem()
	}
	if beanType.Kind() != reflect.Struct { //如果不是struct类型，直接返回
		return fields
	}
	for i := 0; i < beanType.NumField(); i++ {
		field := beanType.Field(i) //通过这个i作为它的索引，从0开始来取得它的字段
		value := beanValue.Field(i)
		if value.IsValid() {
			fields = append(fields, field.Name)
		}
	}
	return fields
}
