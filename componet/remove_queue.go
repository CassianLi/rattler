package componet

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type RemoveParam struct {
	// SourceFile The source file
	SourceFile string

	// MoveTo The directory for move to
	MoveTo string
}

// RemoveQueue 移动文件队列，本质是一个channel
type RemoveQueue struct {
	Queue chan RemoveParam
}

var queue *RemoveQueue

// InitRemoveQueue 初始化队列
func InitRemoveQueue(size int) *RemoveQueue {
	queue = &RemoveQueue{
		Queue: make(chan RemoveParam, size),
	}

	return queue
}

// Add 添加需要移动的文件
func (q *RemoveQueue) Add(param RemoveParam) {
	q.Queue <- param
}

// Run 启动需要的业务，需要传入业务方法
func (q *RemoveQueue) Run() {
	defer close(q.Queue)
	for param := range q.Queue {
		srcFile := param.SourceFile
		filename := filepath.Base(srcFile)
		targetPath := filepath.Join(param.MoveTo, filename)

		log.Infof("%s moving to dir: %s", srcFile, targetPath)
		err := os.Rename(srcFile, targetPath)
		if err != nil {
			log.Errorf("Error moving file, err: %s", err)
		}
	}
}
