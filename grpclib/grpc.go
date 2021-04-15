package grpclib

import (
	"context"
	"strconv"
	"time"

	"github.com/dxasu/gostar/config"
	"github.com/dxasu/gostar/util"

	log "github.com/dxasu/gostar/util/glog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// UserService服务存根，可直接调用服务方法
	userCli UserInfoServiceClient
)

// GrpcInit 初始化grpc, 创建grpc的服务连接
func GrpcInit() {
	addr := config.GetCfgByKey("grpc_config.host")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to connect : ", err)
	}
	userCli = NewUserInfoServiceClient(conn)
}

// uid 100221
func GetUser(uid uint64) *UserInfoResp {
	timeCtx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	token, err := util.EncryptTicket(int(uid))
	if err != nil {
		log.Errorf("Failed to fetch encryted ticket with err(%v)", err)
		return nil
	}

	ctx := metadata.AppendToOutgoingContext(timeCtx, "lang", "cn", "uid", strconv.FormatInt(int64(uid), 10), "x-auth-token", token, "token", token)
	// 测试调用方法Get，返回user对象
	user, err := userCli.GetUserInfo(ctx, &UserInfoReq{Uid: uid})
	if err != nil {
		log.Infof("Get connect failed :%+v", err)
		return nil
	}
	log.Infof("userinfo:%+v", *user)
	return user
}

// func DialTimeout(network, address string, timeout time.Duration) (*rpc.Client, error) {
// 	conn, err := net.DialTimeout(network, address, timeout)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonrpc.NewClient(conn), err
// }

// func TestGetList() {
// 	// 测试调用方法GetList，返回一个Stream流，循环获取多个user对象
// 	recvStream, err := userCli.GetList(context.Background(), &UserReq{Id: 1})
// 	if err != nil {
// 		log.Infof("GetList connect failed :%v", err)
// 		return
// 	}
// 	for {
// 		user, err := recvStream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		log.Println("GetList获取的一条响应数据：", *user)
// 	}
// }

// func TestWaitGet() {
// 	// 测试调用方法WaitGet，传入多条请求数据，返回一个user对象
// 	sendStream, err := userCli.WaitGet(context.Background())
// 	if err != nil {
// 		log.Infof("WaitGet connect failed :%v", err)
// 		return
// 	}
// 	for i := 0; i < 5; i++ { // 一次传入5条请求数据
// 		if err = sendStream.Send(&UserReq{Id: int32(i)}); err != nil {
// 			log.Infof("WaitGet send failed :%v", err)
// 			return
// 		}
// 	}
// 	// 服务端接受全部请求数据后，返回一个user对象
// 	user, err := sendStream.CloseAndRecv()
// 	if err != nil {
// 		log.Infof("WaitGet recv failed :%v", err)
// 		return
// 	}
// 	log.Println("WaitGet响应数据：", *user)
// }
