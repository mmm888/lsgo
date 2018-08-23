package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

func getHardLinkCount(i os.FileInfo) uint16 {
	return i.Sys().(*syscall.Stat_t).Nlink
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

func getFileSize(h bool, i os.FileInfo) string {
	var size string

	if h {
		size = fmt.Sprint(humanateBytes(uint64(i.Size())))
	} else {
		size = fmt.Sprint(i.Size())
	}

	return size
}

func getTimeStamp(i os.FileInfo) string {
	// TODO: 1 年以上前の場合のフォーマット
	return i.ModTime().Format("Jan 02 15:04")
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

func longFormat(o *Options, i os.FileInfo) string {
	fType := i.Mode()
	hardlink := getHardLinkCount(i)
	owner := getUserName(i)
	group := getGroupName(i)
	byteSize := getFileSize(o.human, i)
	timeStamp := getTimeStamp(i)
	name := getFileName(i)

	return fmt.Sprintf("%1s %2d %s %s %5s %s %s", fType, hardlink, owner, group, byteSize, timeStamp, name)
}

func getUsedBlockSize(i os.FileInfo) int {
	return int(i.Sys().(*syscall.Stat_t).Blocks)
}
