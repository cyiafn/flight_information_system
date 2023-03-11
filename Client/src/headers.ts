// | uint8, 1 byte: request type | string, 10 bytes: requestID | uint8: Buffer no. | uint8: no. of Buffer Arrays | rest of payload

import { customAlphabet } from "nanoid";

export const createRequestId = customAlphabet(
  "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
  9
);

// Creating header to send over
export function constructHeaders(
  requestType: number,
  requestIdStr: string,
  byteArrayBufferNo: number,
  totalByteArrayBuffers: number
) {
  const header = Buffer.allocUnsafe(26);

  header.writeUInt8(requestType, 0);
  header.write(requestIdStr, 1);
  header.writeBigInt64LE(BigInt(byteArrayBufferNo), 10);
  header.writeBigInt64LE(BigInt(totalByteArrayBuffers), 18);

  return header;
}

// Destructure header to interpret
export function deconstructHeaders(buffer: Buffer) {
  const bufferSliced = buffer.subarray(0, 26);
  const requestType = bufferSliced.readUInt8(0);
  const requestIdStr = bufferSliced.toString("utf-8", 1, 10);
  const byteArrayBufferNo = bufferSliced.readBigInt64LE(10);
  const totalByteArrayBuffers = bufferSliced.readBigInt64LE(18);

  return {
    requestType: requestType,
    requestId: requestIdStr,
    byteArrayBufferNo: byteArrayBufferNo,
    totalByteArrayBuffers: totalByteArrayBuffers,
  };
}
