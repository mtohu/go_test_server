package routes

import "gotest/common"

type Rpcroutes struct {
	C string
	A string
	Args interface{}
}
type RpcRouter struct {

}

func (r *RpcRouter) RpcAccept(rargs *Rpcroutes,reply *common.ResData) error{
	var data=make(map[string]interface{})
	data["a"]="sss"
	reply.Code=200
	reply.Msg="成功"
	reply.Data=data

	bb,err :=setRoteDispatchs(rargs)
	if(err !=nil){
		common.Log("error dis error",err)
	}
	common.Log("iiiii------:",bb.Data)
	return nil
}
func setRoteDispatchs(param *Rpcroutes)( *common.ResData,error){
	var rr common.ResData
    return &rr,nil
}
