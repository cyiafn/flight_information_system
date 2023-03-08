import dgram from "dgram";
import { constructHeaders } from "./headers";

const server = dgram.createSocket("udp4");

server.on("listening", () => {
  const address = server.address();
  console.log(`Server listening on ${address.address}:${address.port}`);
});

server.on("message", (msg, rinfo) => {
  console.log(`Received message ${rinfo.address}:${rinfo.port}`);
  console.log(msg);

  //sending msg to the client
  // let response = Buffer.from("From server : your msg is received");
  // const header = constructHeaders(107, "abc123", 1, 1);
  let response = Buffer.from([
    0x66, 0x67, 0x54, 0x64, 0x6c, 0x6a, 0x44, 0x30, 0x4d, 0x70, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,
    0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
  ]);
  server.send(response, rinfo.port, "localhost", function (error) {
    if (error) {
      server.close();
    } else {
      console.log(`Data sent to ${rinfo.address}:${rinfo.port}`);
    }
  });
});

const PORT = 8080;

server.bind(PORT, "localhost");
