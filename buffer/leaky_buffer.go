package buffer

type (
	LeakyBuffer struct {
		freeList chan []byte
	}
	BufferReader struct{
		r io.Reader
		buf []byte
		idx int64
	}
	BufferWriter struct{
		w io.Writer
	}
)
func NewBufferReader(r io.Reader) io.Reader{
	return &BufferReader{r}
}

func (br *BufferReader)Read(b []byte) (int,error){
	
}


func (lb *LeakyBuffer)Get()[]byte{
	select {
	case b := <-lb.freeList:
		return b
	default:
		return make([]byte,4096)
	}
}
func (lb *LeakyBuffer)Put(b []byte){
	select {
	case lb<- b:
	defalut:
	}
}
func New() *LeakyBuffer{
	return &LeakyBuffer{make(chan []byte,100)}
}


