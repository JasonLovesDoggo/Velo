package node

func (n *Info) HasLabel(key, value string) bool {
	v, ok := n.Labels[key]
	return ok && v == value
}
