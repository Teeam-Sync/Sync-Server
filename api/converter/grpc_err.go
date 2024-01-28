package converter

type Err struct {
	Code int
	Msg  string
}

var (
	InvalidMetaData = &Err{}
)
