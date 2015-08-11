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
      t;     // temporary
  for (let i = 0; i < s.length;) {
    let c = s.charCodeAt(i++);
    if (c === 61 && j > 1)
        break;
    if ((c = s64[c]) === undefined)
        throw Error(invalidEncoding);
    switch (j) {
      case 0:
        t = c;
        j = 1;
        break;
      case 1:
        buffer[offset++] = t << 2 | (c & 48) >> 4;
        t = c;
        j = 2;
        break;
      case 2:
        buffer[offset++] = (t & 15) << 4 | (c & 60) >> 2;
        t = c;
        j = 3;
        break;
      case 3:
        buffer[offset++] = (t & 3) << 6 | c;
        j = 0;
        break;
    }
  }
  if (j === 1)
      throw Error(invalidEncoding);
  return new Uint8Array(buffer);
}

function b64Test(s: string): boolean {
	return /^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/.test(s);
}

export interface InitReq extends RequestInit {
  pathPrefix?: string
}

export type FetchFn = typeof fetch

let fetchFn: FetchFn = (input, init) => fetch(input, init)

export function setFetchFn(fetcher: FetchFn) {
  fetchFn = fetcher
}

export function resetFetchFn() {
  fetchFn = (input, init) => fetch(input, init)
}

export function replacer(key: any, value: any): any {
  if(value && value.constructor === Uint8Array) {
    return b64Encode(value, 0, value.length);
  }

  return value;
}

export function fetchReq<I, O>(path: string, init?: InitReq): Promise<O> {
  const {pathPrefix, ...req} = init || {}

  const url = pathPrefix ? `${pathPrefix}${path}` : path

  return fetchFn(url, req).then(r => r.json().then((body: O) => {
    if (!r.ok) { throw body; }
    return body;
  })) as Promise<O>
}
