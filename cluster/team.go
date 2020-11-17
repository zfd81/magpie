package cluster

import "sort"

type Team []*Node

func (t *Team) AddMember(node *Node) {
	t.RemoveMember(node.Id) //删除相同节点
	*t = append(*t, node)
}

func (t *Team) RemoveMember(id string) int {
	for i, v := range *t {
		if v.Id == id {
			*t = append((*t)[:i], (*t)[i+1:]...)
			return 1
		}
	}
	return 0
}

func (t *Team) IsMember(node *Node) bool {
	for _, v := range *t {
		if v.Id == node.Id {
			return true
		}
	}
	return false
}

func (t Team) GetLeader() *Node {
	if t.Len() > 0 {
		sort.Stable(t)
		return t[0]
	}
	return nil
}

func (t *Team) IsLeader(node *Node) bool {
	if t.Len() < 1 {
		return true
	}
	l := t.GetLeader()
	if l.Id == node.Id {
		return true
	}
	return false
}

func (t Team) Len() int {
	return len(t)
}
func (t Team) Less(i, j int) bool {
	return t[i].Id < t[j].Id
}
func (t Team) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
