// | uint8, 1 byte: request type | string, 10 bytes: requestID | uint8: Packet no. | uint8: no. of Packets | rest of payload

import { customAlphabet } from "nanoid";

export const createRequestId = customAlphabet(
  "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
  10
);

export function constructHeaders(
  requestType: number,
  requestIdStr: string,
  packetNo: number,
  noOfPackets: number
) {
  const header = Buffer.allocUnsafe(13);

  header.writeUInt8(requestType, 0);
  header.write(requestIdStr, 1);
  header.writeUInt8(packetNo, 11);
  header.writeUInt8(noOfPackets, 12);

  return header;
}

export function deconstructHeaders(packet: Buffer) {
  const packetSliced = packet.subarray(0, 13);
  const requestType = packetSliced.readUint8(0);
  const requestIdStr = packetSliced.toString("utf-8", 1, 11);
  const packetNo = packetSliced.readUInt8(11);
  const noOfPackets = packetSliced.readUInt8(12);

  return {
    requestType: requestType,
    requestId: requestIdStr,
    packetNo: packetNo,
    noOfPackets: noOfPackets,
  };
}
