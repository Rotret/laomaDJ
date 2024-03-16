package service

import (
	"context"
	"regexp"
	"time"

	pb "customer/api/customer"
	"customer/api/verifyCode"
	"customer/internal/data"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	cd *data.CustomerData
}

func NewCustomerService(cd *data.CustomerData) *CustomerService {
	return &CustomerService{
		cd: cd,
	}
}

func (s *CustomerService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCodeRes, error) {
	//一、校验手机号
	pattern := `^(13\d|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18\d|19[0-35-9])\d{8}$`
	regexpPattern := regexp.MustCompile(pattern)   //编译正则表达式字符串
	if !regexpPattern.MatchString(req.Telephone) { //执行匹配操作
		return &pb.GetVerifyCodeRes{
			Code:    1,
			Message: "电话号码格式错误",
		}, nil
	}

	//二、通过验证码生成服务获取验证码（服务间通信，grpc）
	// 连接目标grpc服务器
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("localhost:9000"), //verifyCode grpc service 地址
	)
	if err != nil {
		return &pb.GetVerifyCodeRes{
			Code:    1,
			Message: "验证码服务不可用",
		}, nil
	}
	//关闭
	defer func() {
		_ = conn.Close()
	}()
	//发送验证码请求
	client := verifyCode.NewVerifyCodeClient(conn)
	response, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   1,
	})
	if err != nil {
		return &pb.GetVerifyCodeRes{
			Code:    1,
			Message: "验证码获取失败",
		}, nil
	}

	//三、redis的临时存储
	//连接redis
	const life = 60
	if err:=s.cd.GetVerifyCode(req.Telephone,response.Code,life);err != nil {
		return &pb.GetVerifyCodeRes{
			Code:    1,
			Message: "验证码存储错误（Redis的操作错误）",
		}, nil
	}

	//生成响应
	return &pb.GetVerifyCodeRes{
		Code:           0,
		VerifyCode:     response.Code,
		VerifyCodeTime: time.Now().Unix(),
		VerifyCodeLife: life,
	}, nil
}
