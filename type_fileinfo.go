package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"time"
)

type FileInfo struct {
	fileType      os.FileMode
	hardlinkNum   uint16
	owner         string
	group         string
	byteSize      int64
	timeStamp     time.Time
	fileName      string
	usedBlockSize int64
}

func NewFileInfo(fi os.FileInfo) *FileInfo {
	return &FileInfo{
		fileType:      fi.Mode(),
		hardlinkNum:   fi.Sys().(*syscall.Stat_t).Nlink,
		owner:         getUserName(fi),
		group:         getGroupName(fi),
		byteSize:      fi.Size(),
		timeStamp:     fi.ModTime(),
		fileName:      getFileName(fi),
		usedBlockSize: fi.Sys().(*syscall.Stat_t).Blocks,
	}
}

func (f *FileInfo) LongFormat(o *Options) string {
	return fmt.Sprintf("%1s %2d %s %s %5s %s %s", f.fileType, f.hardlinkNum, f.owner, f.group, f.getFileSize(o.human), f.getTimeStamp(), f.fileName)
}

func (f *FileInfo) GetUsedBlockSize() int {
	return int(f.usedBlockSize)
}

func (f *FileInfo) getFileSize(h bool) string {
	var size string

	if h {
		size = fmt.Sprint(humanateBytes(uint64(f.byteSize)))
	} else {
		size = fmt.Sprint(f.byteSize)
	}

	return size
}

func (f *FileInfo) getTimeStamp() string {
	// TODO: 1 年以上前の場合のフォーマット
	return f.timeStamp.Format("Jan 02 15:04")
}

func getUserName(i os.FileInfo) string {
	uid := i.Sys().(*syscall.Stat_t).Uid
	u, _ := user.LookupId(fmt.Sprint(uid))

	return u.Name
}

func getGroupName(i os.FileInfo) string {
	gid := i.Sys().(*syscall.Stat_t).Gid
	g, _ := user.LookupGroupId(fmt.Sprint(gid))

	return g.Name
}

func getFileName(i os.FileInfo) string {
	var n string

	if i.IsDir() {
		n = fmt.Sprintf("\x1b[36m%s\x1b[0m", i.Name())
	} else {
		n = i.Name()
	}

	return n
}
