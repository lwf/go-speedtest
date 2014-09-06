package transfer

type Worker interface {
	Init(string) chan int
	Sink(int) Sink
	String() string
}
