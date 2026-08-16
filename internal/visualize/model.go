package visualize

type NodeKind string

const (
	Client       NodeKind = "client"
	Proxy        NodeKind = "proxy"
	Application  NodeKind = "application"
	Database     NodeKind = "database"
	Cache        NodeKind = "cache"
	Container    NodeKind = "container"
	Orchestrator NodeKind = "orchestrator"
	Pipeline     NodeKind = "pipeline"
)

type Node struct {
	ID, Label, Description string
	Kind                   NodeKind
}

type Edge struct {
	From, To, Label string
}

type Diagram struct {
	Title string
	Nodes []Node
	Edges []Edge
}
