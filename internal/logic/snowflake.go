package logic

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type snowNode struct {
	node *snowflake.Node
}

func NewIdGenerator(node_id int64, logger *zap.Logger) (*snowNode, error) {
	node, err := snowflake.NewNode(node_id)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create snowflake node")
		return nil, err
	}
	return &snowNode{
		node: node,
	}, nil
}

func (i snowNode) GenID() uint64 {
	id := i.node.Generate().Int64()
	return uint64(id)
}
