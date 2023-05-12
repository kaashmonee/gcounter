package constant

// Server request constants
const (
	display int = iota
	merge
	visit
	value
	update
)

var (
	Request request
)

type request int

func (r request) Display() int { return display }
func (r request) Merge() int   { return merge }
func (r request) Visit() int   { return visit }
func (r request) Value() int   { return value }
func (r request) Update() int  { return update }
