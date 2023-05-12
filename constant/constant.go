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
	ServerRequest serverRequest
	WorkerRequest workerRequest
)

type serverRequest int

func (s serverRequest) Display() int { return display }
func (s serverRequest) Merge() int   { return merge }

// Worker request constants
type workerRequest int

func (w workerRequest) Visit() int  { return visit }
func (w workerRequest) Value() int  { return value }
func (w workerRequest) Update() int { return update }
func (w workerRequest) Merge() int  { return merge }
