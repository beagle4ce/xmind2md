package unarchive

import (
	"archive/zip"
	"fmt"
	errhandle "xmind2md/err-handle"
)

// Unarchived 从ZIP文件中提取content.json内容
// 参数：
//
//	src - 源ZIP文件路径
//
// 返回值：
//
//	提取到的JSON字符串
//
// 功能说明：
//  1. 打开ZIP文件读取器
//  2. 遍历ZIP内文件查找content.json
//  3. 找到后读取并返回其内容
//  4. 未找到则触发panic
func Unarchived(src string) []byte {
	// 打开ZIP文件读取器，使用统一错误处理
	readCloser, err := zip.OpenReader(src)
	errhandle.HandleError(err)

	// 延迟关闭读取器，确保函数退出前释放资源
	defer func() {
		err := readCloser.Close()
		errhandle.HandleError(err)
	}()

	// 遍历ZIP内所有文件
	for _, file := range readCloser.File {
		// 查找目标文件：名为content.json且不是目录
		if file.Name == "content.json" && !file.FileInfo().IsDir() {
			fmt.Println("content.json finded!")
			// 找到后立即读取并返回内容
			return outputJsonString(file)
		}
	}

	// 未找到目标文件时触发异常
	panic("content.json not finded!")
}

// outputJsonString 从ZIP文件条目中读取JSON内容
// 参数：
//
//	file - 目标ZIP文件条目指针
//
// 返回值：
//
//	读取到的完整JSON字符串
//
// 实现特点：
//   - 分块读取（1024字节/次）避免内存过大
//   - 自动处理文件关闭
func outputJsonString(file *zip.File) []byte {
	// 打开文件条目读取器
	readCloser, err := file.Open()
	errhandle.HandleError(err)

	// 延迟关闭确保资源释放
	defer func() {
		err := readCloser.Close()
		errhandle.HandleError(err)
	}()

	// 这里初始len一定要是0, 否则后续会出问题
	jsonByte := make([]byte, 0, 204800)

	// 循环读取直到文件结束
	for {
		// 创建10KB缓冲区（平衡内存效率和IO次数）
		buf := make([]byte, 10240)

		// 读取数据到缓冲区
		end, err := readCloser.Read(buf)

		// 处理读取错误或EOF
		if err != nil {
			break
		}

		// 拼接字节切片
		jsonByte = append(jsonByte, buf[:end]...)
	}
	// 返回读取到的完整JSON字节切片
	return jsonByte
}
