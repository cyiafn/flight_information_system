// | uint8, 1 byte: request type | string, 10 bytes: requestID | uint8: Packet no. | uint8: no. of Packets | rest of payload

import { customAlphabet } from "nanoid";

export const createRequestId = customAlphabet(
  "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
  9
);

// Creating header to send over
export function constructHeaders(
  requestType: number,
  requestIdStr: string,
  packetNo: number,
  noOfPackets: number
) {
  const header = Buffer.allocUnsafe(26);

  header.writeInt8(requestType, 0);
  header.write(requestIdStr, 1);
  header.writeBigInt64LE(BigInt(packetNo), 10);
  header.writeBigInt64LE(BigInt(noOfPackets), 18);

  return header;
}

// Destructure header to interpret
export function deconstructHeaders(packet: Buffer) {
  const packetSliced = packet.subarray(0, 26);
  const requestType = packetSliced.readInt8(0);
  const requestIdStr = packetSliced.toString("utf-8", 1, 10);
  const packetNo = packetSliced.readBigInt64LE(10);
  const noOfPackets = packetSliced.readBigInt64LE(18);

  return {
    requestType: requestType,
    requestId: requestIdStr,
    packetNo: packetNo,
    noOfPackets: noOfPackets,
  };
}
