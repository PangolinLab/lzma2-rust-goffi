package lzma2

/*
	#cgo CFLAGS: -I${SRCDIR}/include
	#cgo LDFLAGS: -lkernel32 -lntdll -luserenv -lws2_32 -ldbghelp -L${SRCDIR}/bin -llzma2
	#include <stdlib.h>
	#include <lzma2_interface.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"unsafe"
)

func init() {
	// 动态库最终路径
	var libFile string
	switch runtime.GOOS {
	case "windows":
		libFile = "bin/lzma2.dll"
	case "darwin":
		libFile = "bin/liblzma2.dylib"
	default:
		libFile = "bin/liblzma2.so"
	}

	// 如果库不存在，则编译 Rust 并复制到 bin/
	if _, err := os.Stat(libFile); os.IsNotExist(err) {
		// Rust 源码目录（Cargo.toml 所在目录）
		rustDir := "../" // 根据你的目录结构调整
		buildCmd := exec.Command("cargo", "build", "--release")
		buildCmd.Dir = rustDir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			panic("Failed to build Rust library: " + err.Error())
		}

		// 源文件路径（默认 target/release/）
		var srcLib string
		switch runtime.GOOS {
		case "windows":
			srcLib = filepath.Join(rustDir, "target", "release", "lzma2.dll")
		case "darwin":
			srcLib = filepath.Join(rustDir, "target", "release", "liblzma2.dylib")
		default:
			srcLib = filepath.Join(rustDir, "target", "release", "liblzma2.so")
		}

		// 确保 bin 目录存在
		_ = os.MkdirAll("bin", 0755)

		// 复制库到 bin/
		input, err := os.ReadFile(srcLib)
		if err != nil {
			panic("Failed to read Rust library: " + err.Error())
		}
		if err := os.WriteFile(libFile, input, 0644); err != nil {
			panic("Failed to write library to bin/: " + err.Error())
		}
	}
}

func Compress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return []byte{}, errors.New("lzma2 compress: empty input")
	}
	cIn := C.CBytes(in)
	defer C.free(cIn)

	var outPtr *C.uint8_t
	var outLen C.size_t

	ret := C.lzma2_compress((*C.uint8_t)(cIn), C.size_t(len(in)),
		&outPtr, &outLen)
	if ret != 0 {
		return nil, errors.New(fmt.Sprint("lzma2 compress failed: ", int(ret)))
	}

	if outPtr == nil {
		return []byte{}, nil
	}

	out := C.GoBytes(unsafe.Pointer(outPtr), C.int(outLen))
	C.lzma2_free(unsafe.Pointer(outPtr))
	return out, nil
}

func Decompress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return []byte{}, errors.New("lzma2 decompress: empty input")
	}

	cIn := C.CBytes(in)
	defer C.free(cIn)

	var outPtr *C.uint8_t
	var outLen C.size_t

	ret := C.lzma2_decompress((*C.uint8_t)(cIn), C.size_t(len(in)),
		&outPtr, &outLen)
	if ret != 0 {
		return nil, errors.New(fmt.Sprint("lzma2 decompress failed: ", int(ret)))
	}

	if outPtr == nil {
		return []byte{}, nil
	}

	out := C.GoBytes(unsafe.Pointer(outPtr), C.int(outLen))
	C.lzma2_free(unsafe.Pointer(outPtr))
	return out, nil
}
