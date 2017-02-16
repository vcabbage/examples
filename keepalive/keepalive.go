package keepalive

/*
#include <stdlib.h>

struct Thing {
    long long *ptr;
    int len;
};

struct Thing* newThing(long long * arr, int len)
{
    struct Thing* t = (struct Thing*)malloc(sizeof(struct Thing));;
    t->ptr = arr;
    t->len = len;
    return t;
}

long long doIt(struct Thing *t)
{
    long long sum = 0;
    for(int i = 0; i < t->len; i++){
        sum += t->ptr[i];
    }
    return sum;
}
*/
import "C"
import (
	"math/rand"
	"runtime"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// DoIt fills a slice of int with values 0 to len(s)-1.
//
// It uses cgo, for reasons.
func DoIt(s []int64) int64 {
	modified := doIt(s)
	return modified
}

// DoItKeepAlive is like DoIt, but uses a call to runtime.KeepAlive.
func DoItKeepAlive(s []int64) int64 {
	modified := doIt(s)
	runtime.KeepAlive(s)
	return modified
}

func doIt(s []int64) int64 {
	cThing := C.newThing(
		// Pointer to first element in the slice's underlying array
		// (passing Go allocated memory to C)
		(*C.longlong)(unsafe.Pointer(&s[0])),
		// Length of the slice
		C.int(len(s)),
	)

	// Make sure we free our cThing
	defer C.free(unsafe.Pointer(cThing))

	// Simulate some work causing GC to run
	runtime.GC()
	other := make([]int, rand.Intn(10000))
	for i := range other {
		other[i] = rand.Int()
	}

	// Return the result of our complex C function
	return int64(C.doIt(cThing))
}
