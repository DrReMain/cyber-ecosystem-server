# 关于接口入参和校验

所有的接口入参字段必须使用 指针类型 *Entity， 包括数组[]*Entity  
go的反序列化会将 `未传`和`null`置为nil，校验的`required`规则是不能为nil
同样`omitempty`时，nil不会被校验  

> [!IMPORTANT]
> 所以前端 request 的数据
> 1. 可以不传   （omitempty情况）
> 2. 不能传null  (required检验不通过，omitempty会nil)
> 3. 可以传0值   (required通过)

# 关于接口的path参数

path参数一定是required的

# 关于接口的query参数

规定数组内元素不能是对象，必须是基本类型

# 关于接口返回值

所有字段也使用 指针类型包括数组  

> [!IMPORTANT]
> 
> 但类型的所有字段都必须明确编写返回，可以是nil -> null  
> 除了数组，没有也要返回空数组
