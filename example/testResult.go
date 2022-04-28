package example

type result struct {
	Code int	`tjson:"code"`
	Tip string	`tjson:"tip"`
	Data interface{}	`tjson:"data"`
}

//	创建一个Result实例
//	tip string	提示
//	data interface{}	数据
func resultNewSuccess(tip string,data interface{}) result{
	if "" == tip {
		tip = "成功"
	}
	return result{2000,tip,data}
}


