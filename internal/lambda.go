// Generated by `wit-bindgen` 0.17.0. DO NOT EDIT!
package internal

// #include "lambda.h"
import "C"
import "unsafe"

// Import functions from manetu:lambda/sparql@0.0.2
func ManetuLambda0_0_2_SparqlQuery(expr string) string {
	var lower_expr C.lambda_string_t

	// use unsafe.Pointer to avoid copy
	lower_expr.ptr = (*uint8)(unsafe.Pointer(C.CString(expr)))
	lower_expr.len = C.size_t(len(expr))
	var ret C.lambda_string_t
	C.manetu_lambda_sparql_query(&lower_expr, &ret)
	var lift_ret string
	lift_ret = C.GoStringN((*C.char)(unsafe.Pointer(ret.ptr)), C.int(ret.len))
	return lift_ret
}

// Export functions from manetu:lambda/guest@0.0.2
var exports_manetu_lambda0_0_2_guest ExportsManetuLambda0_0_2_Guest = nil

// `SetExportsManetuLambda0_0_2_Guest` sets the `ExportsManetuLambda0_0_2_Guest` interface implementation.
// This function will need to be called by the init() function from the guest application.
// It is expected to pass a guest implementation of the `ExportsManetuLambda0_0_2_Guest` interface.
func SetExportsManetuLambda0_0_2_Guest(i ExportsManetuLambda0_0_2_Guest) {
	exports_manetu_lambda0_0_2_guest = i
}

type ExportsManetuLambda0_0_2_Guest interface {
	HandleRequest(request string) string
	Malloc(len uint32) uint32
	Free(ptr uint32)
}

//export exports_manetu_lambda_guest_handle_request
func exportsManetuLambda002GuestHandleRequest(request *C.lambda_string_t, ret *C.lambda_string_t) {
	var lift_request string
	lift_request = C.GoStringN((*C.char)(unsafe.Pointer(request.ptr)), C.int(request.len))
	result := exports_manetu_lambda0_0_2_guest.HandleRequest(lift_request)
	var lower_result C.lambda_string_t

	// use unsafe.Pointer to avoid copy
	lower_result.ptr = (*uint8)(unsafe.Pointer(C.CString(result)))
	lower_result.len = C.size_t(len(result))
	*ret = lower_result

}

//export exports_manetu_lambda_guest_malloc
func exportsManetuLambda002GuestMalloc(len C.uint32_t) C.uint32_t {
	var lift_len uint32
	lift_len = uint32(len)
	result := exports_manetu_lambda0_0_2_guest.Malloc(lift_len)
	lower_result := C.uint32_t(result)
	return lower_result

}

//export exports_manetu_lambda_guest_free
func exportsManetuLambda002GuestFree(ptr C.uint32_t) {
	var lift_ptr uint32
	lift_ptr = uint32(ptr)
	exports_manetu_lambda0_0_2_guest.Free(lift_ptr)

}
