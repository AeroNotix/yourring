package main

/*
#cgo LDFLAGS: -luring
#include <liburing.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

type IOData struct {
	read                C.int
	firstOffset, offset C.off_t
	firstLen            C.size_t
}

type Ring struct {
	ring C.struct_io_uring
}

func (r *Ring) QueueRead(fd uintptr, size uint64, offset uint64) {

	data := IOData{}
	sqe := C.io_uring_get_sqe(&r.ring)
	if sqe == nil {
		return
	}

	data.read = 1
	data.offset = C.long(offset)
	data.firstOffset = C.long(offset)
	data.firstLen = C.ulong(size)

	iov := (*C.struct_iovec)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_iovec{}))))
	if iov == nil {
		panic("LOL")
	}

	iov.iov_base = unsafe.Pointer(&data.firstOffset)
	iov.iov_len = C.ulong(size)

	C.io_uring_prep_readv(sqe, C.int(fd), iov, 1, C.long(offset))
	C.io_uring_sqe_set_data(sqe, unsafe.Pointer(iov))
	C.io_uring_submit(&r.ring)

	// Something like this should be used: figure out how to do it.
	// var cqe []C.struct_io_uring_cqe
	// if C.io_uring_wait_cqe(&r.ring, &cqe) < 0 {
	// 	panic("lolno")
	// }

	// just print the iovec to stdout, works because the iovec is
	// what's populated by the actual io_yourring shit.
	//
	// the iovec sometimes has no data, because async - run it a few
	// times.
	C.writev(1, iov, 3)
}

func (r *Ring) Init() error {
	if ret := C.io_uring_queue_init(512, &r.ring, 0); ret < 0 {
		return fmt.Errorf("cannot initialise ring: %d", ret)
	}
	return nil
}

func main() {
	r := &Ring{}
	fmt.Println(r.Init())
	f, err := os.Open("/dev/urandom")
	if err != nil {
		fmt.Println(err)
	}
	r.QueueRead(f.Fd(), 1024, 0)
}
