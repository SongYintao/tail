package watch

type FileChanges struct {
	Modified chan bool
	Truncated chan bool
	Deleted chan bool
}

func NewFileChanges()*FileChanges{
	return &FileChanges{
		Modified:  make(chan bool,1),
		Truncated: make(chan bool,1),
		Deleted:   make(chan bool,1),
	}
}

func (fc *FileChanges) NotifyModified() {
	sendOnlyIfEmpty(fc.Modified)
}

func (fc *FileChanges) NotifyTruncated() {
	sendOnlyIfEmpty(fc.Truncated)
}

func (fc *FileChanges) NotifyDeleted() {
	sendOnlyIfEmpty(fc.Deleted)
}

// sendOnlyIfEmpty sends on a bool channel only if the channel has no
// backlog to be read by other goroutines. This concurrency pattern
// can be used to notify other goroutines if and only if they are
// looking for it (i.e., subsequent notifications can be compressed
// into one).
func sendOnlyIfEmpty(ch chan bool) {
	select {
	case ch <- true:
	default:
	}
}
