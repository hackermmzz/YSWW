package main

import (
	"reflect"
	"fmt"
)

type StructConstraint interface {
    ~struct{} // 约束为结构体类型（包括匿名结构体）
}

type DataBaseOp[T StructConstraint] struct{
	table string
};


func NewDataBaseOp[T StructConstraint](table string)* DataBaseOp[T]{
	ret:=new(DataBaseOp[T]);
	ret.table=table
	return ret
}

func (dbop*DataBaseOp[T])Select(key string,keyValue string)([]T,error){
	ret:=make([]T,0)
	//
	members:=GetAllMember(*new(T))
	find:=MergeByCommaAndQuo(members...)
	cmd:="select "+find+" from "+dbop.table+" where "+key+"="+keyValue

	row,err:= db.Query(cmd)
	defer row.Close()
	if err!=nil{
		return ret,err
	} 
	
	//获取T的类型信息
	tp:=reflect.TypeOf(new(T)).Elem()

	for row.Next(){
		//
		instance:=reflect.New(tp).Elem()
		ptr:=make([]interface{},tp.NumField())
		for i:=0;i<tp.NumField();i+=1{
			ptr[i]=instance.Field(i).Addr().Interface()
		}
		//
		err = row.Scan(ptr...)
		if err != nil {
			return ret, err
		}
		ret = append(ret,instance.Interface().(T))
	}
	//
	return ret,nil
}

func (dbop*DataBaseOp[T])Insert(keyValue string,info T)error{
	value:=keyValue
	//获取info里面所有的值
	t:=reflect.TypeOf(info)
	v:=reflect.ValueOf(info)
	for i:=0;i<t.NumField();i+=1{
		//如果keyvalue为0表示根本不用keyvalue
		if len(value)!=0{
			value+=","
		}
		field0:=t.Field(i)
		field1:=v.Field(i)
		if field0.Type.String()=="string"{
			value+="\""+field1.String()+"\""
		}else{
			value+=fmt.Sprintf("%v",field1.Interface())
		}
	}
	//
	cmd:="insert into " + dbop.table + " values(" +value+");"
	_, err:= db.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}
func GetAllMember(info interface{})[]string {
	ret0 := make([]string, 0)
	//
	t := reflect.TypeOf(info)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i += 1 {
			ret0 = append(ret0, t.Field(i).Name)
		}
	} else {
		ret0 = append(ret0, t.Name())
	}
	//
	return ret0
}

func MarkByQuo(str string)string{
	return "\""+str+"\"";
}

func MergeByCommaAndQuo(str...string)string{
	ans:=""
	flag:=true
	for _,s:=range str{
		if flag{
			flag=false
		}else{
			ans+=","
		}
		ans+=MarkByQuo(s)
	}
	return ans
}


//