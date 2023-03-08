import dgram from "dgram";
import { constructHeaders, deconstructHeaders } from "./headers";
import { RequestType } from "./interfaces";

const server = dgram.createSocket("udp4");

const PORT = 8080;

server.bind(PORT, "127.0.0.1");

server.on("listening", () => {
  const address = server.address();
  console.log(`Server listening on ${address.address}:${address.port}`);
});

server.on("message", (msg, rinfo) => {
  console.log(`Received message ${rinfo.address}:${rinfo.port}`);
  console.log(msg);

  const headers = deconstructHeaders(msg);
  let packet = Buffer.from([0x00]),
    header = Buffer.from([0x00]),
    payload = Buffer.from([0x00]);
  if (headers.requestType === RequestType.GetFlightIdentifiersRequestType) {
    header = constructHeaders(102, "abcdefg", 1, 1);
    payload = Buffer.alloc(21);
    payload.writeInt8(1);
    payload.writeBigInt64LE(BigInt(3), 1);
    payload.writeInt32LE(1, 9);
    payload.writeInt32LE(2, 13);
    payload.writeInt32LE(3, 17);
  } else if (
    headers.requestType === RequestType.GetFlightInformationRequestType
  ) {
    header = constructHeaders(103, "abcdefg", 1, 1);
    payload = Buffer.alloc(21);
    payload.writeInt8(1);
    payload.writeBigInt64LE(BigInt(1678283107), 1);
    payload.writeDoubleLE(200.15, 9);
    payload.writeInt32LE(56, 17);
  } else if (
    headers.requestType === RequestType.MakeSeatReservationRequestType
  ) {
    header = constructHeaders(104, "abcdefg", 1, 1);
    payload = Buffer.alloc(21);
    payload.writeInt8(1);
  } else if (headers.requestType === RequestType.UpdateFlightPriceRequestType) {
    header = constructHeaders(106, "abcdefg", 1, 1);
    payload = Buffer.alloc(31);
    payload.writeInt8(1);
    payload.writeInt32LE(40, 1);
    payload.write("sg\0", 5);
    payload.write("my\0", 8);
    payload.writeBigInt64LE(BigInt(1678283107), 11);
    payload.writeDoubleLE(142.5, 19);
    payload.writeInt32LE(25, 27);
  } else if (headers.requestType === RequestType.CreateFlightRequestType) {
    header = constructHeaders(107, "abcdefg", 1, 1);
    payload = Buffer.alloc(5);
    payload.writeInt8(1);
    payload.writeInt32LE(25, 1);
  }

  packet = Buffer.concat([header, payload], header.length + payload.length);
  server.send(packet, rinfo.port, "localhost", (err) => {
    console.log(packet);
  });
});
