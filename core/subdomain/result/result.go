package result
//final result
type Result struct {
	subdomain string
	record    string
	code      int
}

var FinalResults = make(chan *Result,10)

func (r *Result) GetSubdomain() string  { return r.subdomain }
func (r *Result) GetRecord() string     { return r.record }
func (r *Result) GetCode() int          { return r.code }
func (r *Result) SetSubdomain(v string) { r.subdomain = v }
func (r *Result) SetRecord(v string)    { r.record = v }
func (r *Result) SetCode(v int)         { r.code = v }
