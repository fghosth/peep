package ast

type NodeType int

const (
	NArg    NodeType = iota
	NString          //string 节点
	NIf
	NTrim
	NForEach
	NChoose
	NOtherwise
	NWhen
	NBind
	NInclude
	NWhere
)

func (it NodeType) ToString() string {
	switch it {
	case NString:
		return "NString"
	case NIf:
		return "NIf"
	case NTrim:
		return "NTrim"
	case NForEach:
		return "NForEach"
	case NChoose:
		return "NChoose"
	case NOtherwise:
		return "NOtherwise"
	case NWhen:
		return "NWhen"
	case NBind:
		return "NBind"
	case NInclude:
		return "NInclude"
	case NWhere:
		return "NWhere"
	}
	return "Unknow"
}
