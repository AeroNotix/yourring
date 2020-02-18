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
	read                int
	firstOffset, offset C.off_t
	firstLen            C.size_t
	iov                 C.struct_iovec
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

	data.iov.iov_base = unsafe.Pointer(&data.read)
	data.iov.iov_len = C.ulong(size)

	C.io_uring_prep_readv(sqe, C.int(fd), &data.iov, 1, C.long(offset))
	C.io_uring_sqe_set_data(sqe, unsafe.Pointer(&data.iov))
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
	f, err := os.Open("/dev/null")
	if err != nil {
		fmt.Println(err)
	}
	r.QueueRead(f.Fd(), 1024, 0)
}
