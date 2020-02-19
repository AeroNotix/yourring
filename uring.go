package yourring

/*
#cgo LDFLAGS: -luring
#include <liburing.h>
*/
import "C"
import (
	"fmt"
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
		panic("malloc nullptr")
	}

	iov.iov_base = unsafe.Pointer(&data.firstOffset)
	iov.iov_len = C.ulong(size)

	C.io_uring_prep_readv(sqe, C.int(fd), iov, 1, C.long(offset))
	C.io_uring_sqe_set_data(sqe, unsafe.Pointer(iov))
	C.io_uring_submit(&r.ring)

	cqe := C.struct_io_uring_cqe{}
	cqes := (**C.struct_io_uring_cqe)(unsafe.Pointer(unsafe.Pointer(&cqe)))
	if C.io_uring_wait_cqe(&r.ring, cqes) < 0 {
		panic("cqe read")
	}
	cqePtr := (*C.struct_io_uring_cqe)(unsafe.Pointer(&cqe))
	readData := (*C.struct_iovec)(C.io_uring_cqe_get_data(cqePtr))
	if cqePtr.res < 0 {
		panic("cqe res")
	}
	C.io_uring_cqe_seen(&r.ring, cqePtr)
	fmt.Println(*(*C.char)(readData.iov_base))
}

func (r *Ring) Init() error {
	if ret := C.io_uring_queue_init(512, &r.ring, 0); ret < 0 {
		return fmt.Errorf("cannot initialise ring: %d", ret)
	}
	return nil
}
