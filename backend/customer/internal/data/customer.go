package data

import "time"

//customer 中与数据操作相关的代码
type CustomerData struct{
	data *Data
}

//New 方法
func NewCustomerData(data *Data) *CustomerData {
	return &CustomerData{data: data}
}

//数据操作函数
func (cd *CustomerData) GetVerifyCode (telephone ,code string,ex int64) error{
	//设置key，customer-verify-code
	status := cd.data.Rdb.Set("CVC:"+telephone, code, time.Duration(ex)*time.Second)
	if _, err := status.Result(); err != nil {
		return err
	}
	return nil
}