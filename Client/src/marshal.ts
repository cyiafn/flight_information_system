import { Buffer } from "buffer";

function toByteArray<T>(data: T): Buffer {
  console.log(typeof data);
  if (Buffer.isBuffer(data)) return data;
  else if (typeof data === "string") {
    return Buffer.from(data + "\0");
  } else if (typeof data === "number") {
    var buffer;
    if (!Number.isInteger(data)) {
      buffer = Buffer.alloc(8);
      buffer.writeDoubleLE(data);
    } else {
      buffer = Buffer.alloc(4);
      buffer.writeInt32LE(data);
    }
    return buffer;
  } else if (typeof data === "bigint") {
    buffer = Buffer.alloc(8);
    buffer.writeBigInt64LE(BigInt(data));
    return buffer;
  } else if (typeof data === "boolean") {
    const value = data ? 1 : 0;
    return Buffer.from([value]);
  } else if (Array.isArray(data)) {
    const buffer = Buffer.alloc(1);
    buffer.writeBigInt64LE(BigInt(data.length));
    const bytes = data.map((e) => toByteArray(e));
    return Buffer.concat([buffer, ...bytes]);
  } else if (data instanceof Object) {
    const values = Object.values(data);
    const bytes = values.map((v) => toByteArray(v));
    return Buffer.concat(bytes);
  } else {
    throw new Error(`Unsupported data type: ${typeof data}`);
  }
}

export function marshal(data: any): any {
  let requestBuffer = toByteArray(data);
  console.log(requestBuffer);
  return requestBuffer;
}
