package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	errhandle "xmind2md/err-handle"
	analysis "xmind2md/json-analysis"
	"xmind2md/markdown"
	"xmind2md/unarchive"
)

func main() {
	// 获取当前程序运行路径
	// var runtimeLocation string
	// 获取输入文件路径参数
	var filePathStr string

	/*
		输出:
			0 Executed location: ./compile/xmind2md
			1 First Argument: 123
	*/
	allArgs := os.Args

	// 若输入参数数组小于2,则输出提示信息并退出程序
	if len(allArgs) < 2 {
		fmt.Println("Usage: xmind2md [filePath]")
		return
	}

	// 参数数组大于等于2时, 遍历数组, 获取第1和第2个参数
	for idx, arg := range allArgs {
		switch idx {
		case 0:
			// runtimeLocation = arg
		case 1:
			filePathStr = arg
		}
	}

	// 解压缩文件
	jsonBytes := unarchive.Unarchived(filePathStr)
	fmt.Println("Analysis Json")
	sheets := analysis.JsonAnalysis(jsonBytes)
	fmt.Println("Trans to Markdown")
	// 解析Sheet对象集合,转为MarkDown 字符切片
	mdMap := markdown.ToMarkDown(sheets)
	fmt.Println("Create MarkDown files and Writed")
	// 将每个map的值单独压入新的md文件中
	createMDFiles(filePathStr, mdMap)
	fmt.Println("Xmind to MarkDown files was successed!")
}

// 根据文件路径和md切片映射表生成md文件
func createMDFiles(filePathStr string, mdMap map[string][]string) {
	// 截取文件路径,不包含文件本体
	normalizedPath := filepath.Clean(filePathStr)
	filePath := filepath.Dir(normalizedPath)
	// 循环mdMap结构, 生成对应的md文件
	for key, val := range mdMap {
		// 拼接文件路径和文件名
		fullFilePath := filePath + string(filepath.Separator) + key + ".md"
		// 检查文件是否存在, 存在则删除
		checkFileExist(fullFilePath)
		// 创建文件
		file, err := os.Create(fullFilePath)
		errhandle.HandleError(err)
		// 设置文件权限为0766
		err = file.Chmod(0766)
		errhandle.HandleError(err)
		// 先检查错误, 再设置自动关闭, 否则错误后, 无实例可关
		defer file.Close()
		// 循环写入md切片
		bufWriteMdSlice(file, val)
	}
}

// 5k缓冲写入md切片
func bufWriteMdSlice(file *os.File, vals []string) {
	// 创建缓冲写入器, 5k缓冲
	writer := bufio.NewWriterSize(file, 5*1024)
	// 强制刷新剩余缓存
	defer writer.Flush()
	for _, val := range vals {
		_, err := writer.WriteString(val)
		errhandle.HandleError(err)
	}
}

// 检查文件是否存在, 存在则删除
func checkFileExist(filePath string) {
	// 判断文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 删除已存在的文件 [^1]
		if err := os.Remove(filePath); err != nil {
			errhandle.HandleError(err)
		}
	} else if !os.IsNotExist(err) { // 非不存在的其他错误类型
		errhandle.HandleError(err) // 权限不足等异常场景
	}
}
