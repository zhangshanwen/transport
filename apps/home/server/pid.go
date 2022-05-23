package server

import (
	"bufio"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhangshanwen/transport/common"
	"os"
	"strconv"
)

func (t *Transponder) WritePid2File() (err error) {
	_, err = os.Stat(common.HomePid)
	if err == nil {
		return errors.New("application is running")
	}

	var file *os.File
	if file, err = os.OpenFile(common.HomePid, os.O_WRONLY|os.O_CREATE, os.ModePerm); err != nil {
		return
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	if _, err = write.WriteString(strconv.Itoa(os.Getpid())); err != nil {
		return
	}
	//Flush将缓存的文件真正写入到文件中
	return write.Flush()
}

func (t *Transponder) RemovePid() {
	var err error
	if err = os.Remove(common.HomePid); err != nil {
		logrus.Error("pid文件删除失败 ", err)
	} else {
		logrus.Info("pid文件删除成功 ")

	}
	return
}
