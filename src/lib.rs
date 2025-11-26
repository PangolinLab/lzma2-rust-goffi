use std::io::{Read, Write};
use std::slice;
use std::ptr;
use std::ffi::c_void;

use libc::{size_t, c_int};
use xz2::read::XzDecoder;
use xz2::write::XzEncoder;

const ERR_NULL_PTR: c_int = 1;
const ERR_IO: c_int = 2;
const ERR_ALLOC: c_int = 3;

#[no_mangle]
pub extern "C" fn lzma2_compress(
    input_ptr: *const u8,
    input_len: size_t,
    out_ptr: *mut *mut u8,
    out_len: *mut size_t,
) -> c_int {
    if input_ptr.is_null() || out_ptr.is_null() || out_len.is_null() {
        return ERR_NULL_PTR;
    }

    let input = unsafe { slice::from_raw_parts(input_ptr, input_len as usize) };

    let mut encoder = XzEncoder::new(Vec::new(), 6);
    if encoder.write_all(input).is_err() {
        return ERR_IO;
    }

    let compressed = match encoder.finish() {
        Ok(b) => b,
        Err(_) => return ERR_IO,
    };

    let out_size = compressed.len();
    if out_size == 0 {
        unsafe {
            *out_ptr = std::ptr::null_mut();
            *out_len = 0;
        }
        return 0;
    }

    unsafe {
        let mem = libc::malloc(out_size) as *mut u8;
        if mem.is_null() {
            return ERR_ALLOC;
        }
        ptr::copy_nonoverlapping(compressed.as_ptr(), mem, out_size);
        *out_ptr = mem;
        *out_len = out_size;
    }

    0
}

#[no_mangle]
pub extern "C" fn lzma2_decompress(
    input_ptr: *const u8,
    input_len: size_t,
    out_ptr: *mut *mut u8,
    out_len: *mut size_t,
) -> c_int {
    if input_ptr.is_null() || out_ptr.is_null() || out_len.is_null() {
        return ERR_NULL_PTR;
    }

    let input = unsafe { slice::from_raw_parts(input_ptr, input_len as usize) };
    let mut decoder = XzDecoder::new(input);

    let mut out = Vec::new();
    if decoder.read_to_end(&mut out).is_err() {
        return ERR_IO;
    }

    let out_size = out.len();
    if out_size == 0 {
        unsafe {
            *out_ptr = std::ptr::null_mut();
            *out_len = 0;
        }
        return 0;
    }

    unsafe {
        let mem = libc::malloc(out_size) as *mut u8;
        if mem.is_null() {
            return ERR_ALLOC;
        }
        ptr::copy_nonoverlapping(out.as_ptr(), mem, out_size);
        *out_ptr = mem;
        *out_len = out_size;
    }

    0
}

#[no_mangle]
pub extern "C" fn lzma2_free(ptr: *mut c_void) {
    if ptr.is_null() {
        return;
    }
    unsafe { libc::free(ptr) };
}
