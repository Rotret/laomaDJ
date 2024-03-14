package service

import (
	"context"
	"math/rand"
	"strings"

	pb "verify-code/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode(int(req.Length), req.Type),
	}, nil
}

// 区分验证码类型
func RandCode(l int, t pb.TYPE) string {
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough //与下行结果一样
	case pb.TYPE_DIGIT:
		return randCode("0123456789", l,4)
	case pb.TYPE_LETTER:
		return randCode("abcdefghijklmnopqrstuvwxyz", l,5)
	case pb.TYPE_MIXED:
		return randCode("0123456789abcdefghijklmnopqrstuvwxyz", l,6)
	default:

	}
	return ""
}

// 根据验证码类型生成随机数
func randCode(chars string, l int,idxBits int) string {
	//形成掩码，例：00000000 00000000 00000000 00111111
	idxMask:=1<<idxBits-1
	//63位最大可以使用次数
	idxMax:=63/idxBits

	//利用string builder构建结果缓冲
	sb:=strings.Builder{}
	sb.Grow(1)//预定义容量

	//i 索引，cache 随机数缓存，remain 随机数还能用几次
	for i,cache,remain := 0,rand.Int63(),idxMax;i<l;{
		//随机数缓存不足，重新生成
		if remain==0 {
			cache,remain = rand.Int63(),idxMax
		}

		//利用掩码生成随机索引，有效索引需小于字符集长度
		//& 与运算：遇0得0，遇1得本身，得到后几位二进制
		if idx:=int(cache & int64(idxMask));idx < len(chars){
			sb.WriteByte(chars[idx])
			i++
		}

		//使用下一组随机数位
		cache>>=idxBits
		remain--
	}		
	
	return sb.String()
}
