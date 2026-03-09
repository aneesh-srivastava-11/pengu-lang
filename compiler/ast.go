package compiler

type Action struct {
	Type string   // "log" or "respond"
	Args []string // arguments for the action
	Line int      // for error reporting
}

type Route struct {
	Method  string
	Path    string
	Actions []Action
	Line    int
}

type Service struct {
	Version        string
	Name           string
	Routes         []Route
	Middleware     []string
	HealthEnabled  bool
	MetricsEnabled bool
	HasJson        bool
	HasAuth        bool
	Line           int
}
