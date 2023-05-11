package utils
import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"project/group_chat/common/message"

	//
	//"hot100/chatsystem2/common/message"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf [8096]byte //这时传输时，使用缓冲
}


func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)
	fmt.Println("读取服务器发送的数据...")
	//fmt.Println("请选择(1-4):")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了 conn 则，就不会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4] 转成一个 uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	//根据 pkgLen 读取消息内容
	n, err1 := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err1 != nil {
		//err = errors.New("read pkg body error")
		return 
	}
	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes！！
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return 
	}
	return 
}


func (this *Transfer) WritePkg(data []byte) (err error) {

	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data)) 
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err2 := this.Conn.Write(this.Buf[:4])
	if n != 4 || err2 != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return 
	}

	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return 
	}
	return 
}