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
  server.on("message", (message, rinfo) => {
    console.log(message);
  });
});

const PORT = 8080;

server.bind(PORT, "localhost");
