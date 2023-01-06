/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

/**
 * base64 encoder and decoder
 * Copied and adapted from https://github.com/protobufjs/protobuf.js/blob/master/lib/base64/index.js
 */
// Base64 encoding table
const b64 = new Array(64);

// Base64 decoding table
const s64 = new Array(123);

// 65..90, 97..122, 48..57, 43, 47
for (let i = 0; i < 64;)
    s64[b64[i] = i < 26 ? i + 65 : i < 52 ? i + 71 : i < 62 ? i - 4 : i - 59 | 43] = i++;

export function b64Encode(buffer: Uint8Array, start: number, end: number): string {
	let parts: string[] = null;
  const chunk = [];
  let i = 0, // output index
    j = 0, // goto index
    t;     // temporary
  while (start < end) {
    const b = buffer[start++];
    switch (j) {
      case 0:
        chunk[i++] = b64[b >> 2];
        t = (b & 3) << 4;
        j = 1;
        break;
      case 1:
        chunk[i++] = b64[t | b >> 4];
        t = (b & 15) << 2;
        j = 2;
        break;
      case 2:
        chunk[i++] = b64[t | b >> 6];
        chunk[i++] = b64[b & 63];
        j = 0;
        break;
    }
    if (i > 8191) {
      (parts || (parts = [])).push(String.fromCharCode.apply(String, chunk));
      i = 0;
    }
  }
  if (j) {
    chunk[i++] = b64[t];
    chunk[i++] = 61;
    if (j === 1)
      chunk[i++] = 61;
  }
  if (parts) {
    if (i)
      parts.push(String.fromCharCode.apply(String, chunk.slice(0, i)));
    return parts.join("");
  }
  return String.fromCharCode.apply(String, chunk.slice(0, i));
}

const invalidEncoding = "invalid encoding";

export function b64Decode(s: string): Uint8Array {
	const buffer = [];
	let offset = 0;
  let j = 0, // goto index
