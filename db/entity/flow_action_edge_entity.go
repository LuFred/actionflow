package entity

type FlowActionEdgeEntity struct {
	Id            int64  `db:"id" description:""`
	EntryEdgeId   int64  `db:"entryEdgeId" description:"作为这条隐含边的创建原因的起始顶点的传入边的Id; 直接边包含与Id列相同的值"`
	DirectEdgeId  int64  `db:"directEdgeId" description:"导致创建此隐含边的直接边的 Id; 直接边包含与Id列相同的值"`
	ExitEdgeId    int64  `db:"exitEdgeId" description:"作为这条隐含边的创建原因的结束顶点的出边的Id;直接边包含与Id列相同的值"`
	StartActionId string `db:"startActionId" description:"开始节点"`
	EndActionId   string `db:"endActionId" description:"结束节点"`
	Hops          int32  `db:"hops" description:"跳数"`
}

func (d *FlowActionEdgeEntity) TableName() string {
	return "flow_action_edge"
}
