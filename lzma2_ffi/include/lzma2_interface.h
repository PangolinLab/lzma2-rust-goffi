#ifndef LZMA2_INTERFACE_H
#define LZMA2_INTERFACE_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

int lzma2_compress(const uint8_t *input_ptr, size_t input_len,
                   uint8_t **out_ptr, size_t *out_len);

int lzma2_decompress(const uint8_t *input_ptr, size_t input_len,
                     uint8_t **out_ptr, size_t *out_len);

void lzma2_free(void *ptr);

#ifdef __cplusplus
}
#endif

#endif // LZMA2_INTERFACE_H
