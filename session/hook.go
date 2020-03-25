package session

import (
	"reflect"
	"zorm/log"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)
//
//func CallMethod(s Session,method string, value interface{}) Session  {
//	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
//	if value != nil {
//		fm = reflect.ValueOf(value).MethodByName(method)
//	}
//	param := []reflect.Value{reflect.ValueOf(s)}
//	if fm.IsValid() {
//		if v := fm.Call(param); len(v) > 0 {
//			if err, ok := v[0].Interface().(error); ok {
//				log.Error(err)
//			}
//		}
//	}
//	return
//}