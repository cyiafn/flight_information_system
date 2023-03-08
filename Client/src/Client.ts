import dgram, { RemoteInfo } from "dgram";
import { Buffer } from "buffer";
import { marshal } from "./marshal";
import {
  constructHeaders,
  createRequestId,
  deconstructHeaders,
} from "./headers";
import { PendingRequest, RequestType } from "./interfaces";
import { unmarshal } from "./unmarshal";
import { exit } from "process";
import { request } from "http";

// Connect to Server Via UDP Connection
// port = 8080

export class UDPClient {
  address: string;
  sendPort: number;
  receivePort: number;
  client: dgram.Socket;
  requestId: string;
  pendingRequests: Map<string, PendingRequest>;
  timeout: number;
  monitorTimeOut: number;

  constructor(
    address: string,
    sendPort: number,
    receivePort: number // Specify the port number (agreed upon) and local ip address of the server.
  ) {
    this.address = address; //IP Address of Server
    this.sendPort = sendPort; //Sending Port Number of Client
    this.receivePort = receivePort; //Listening Port of Client
    this.client = dgram.createSocket("udp4");
    this.requestId = ""; //Tracking of Request ID
    this.pendingRequests = new Map(); //Storing of Pending Requests
    this.timeout = 5000; //Time out when no reply in 0.5 seconds
    this.monitorTimeOut = 0;

    // this.client.bind(this.receivePort, this.address);

    // this.client.on('listening', () => {
    //     console.log(`Test Client listening on ${this.address}:${this.receivePort}`);
    // });

    // this.client.on('message', () => {
    //     console.log(`Test Client received message from ${this.address}:${this.receivePort}`);
    // });
  }

  //UDPClient get Methods
  setRequestId(id: string) {
    this.requestId = id;
  }

  setPendingRequests(id: string, bufferData: string) {
    let attempt = 0;
    if (this.pendingRequests.has(id))
      attempt = (this.pendingRequests.get(id)?.attempts || 0) + 1;

    this.pendingRequests.set(id, { data: bufferData, attempts: attempt });
  }

  getPendingRequests(id: string) {
    return this.pendingRequests.get(id);
  }

  public sendMultipleRequests(dto: any, requestType: number) {
    // If dto is q4
    if (requestType === 4)
      this.monitorTimeOut = dto.LengthOfMonitorIntervalInSeconds;

    // 26 bytes header message max 499 Bytes

    // Create Request Id for these spliced msg requests
    const requestId = createRequestId();
    this.setRequestId(requestId);
    // Marshal data
    const payload = marshal(dto);

    // If data is more than 499 Bytes, send in another packet
    let lengthPayload = payload.length;
    const totalPackets = Math.ceil(lengthPayload / 486);

    for (let i = 1; i <= totalPackets; i++) {
      this.sendRequest(payload, requestType, i, totalPackets);
    }
  }

  private receiveResponse(buffer: Buffer) {
    const header = deconstructHeaders(buffer);
    let tempBuffer;
    while (header.noOfPackets !== header.packetNo) {
      //continue to listen
    }
    const payload = unmarshal(buffer.subarray(26, 512), header.requestType);
    // display the payload information here
    console.log(payload);
  }

  private callback(expireTime: number) {
    console.log("In callback");
    setTimeout(() => {
      console.log("No more monitoring");
      this.client.close();
    }, expireTime * 1000);

    this.client.on("message", (msg) => {
      // unmarshal message
      this.receiveResponse(msg);
    });
  }

  //Send Method to cover both Idempotent and Non-Idempotent
  public sendRequest(
    payload: Buffer,
    requestType: number,
    packetNo: number,
    noOfPackets: number
  ) {
    const header = constructHeaders(
      requestType,
      this.requestId,
      packetNo,
      noOfPackets
    );
    // Converts the message object into array
    const packet = Buffer.concat([header, payload], 512);
    this.client.send(
      packet,
      0,
      packet.length,
      this.sendPort,
      this.address,
      (err: Error | null, bytes: number) => {
        if (err) {
          console.log(`Error sending message: ${err}`);
        } else {
          console.log(`Sent ${bytes} bytes to server`);

          if (requestType === 4) this.callback(this.monitorTimeOut);
          else {
            this.client.on("message", (msg) => {
              clearTimeout(closeSocketTimeout);
              clearTimeout(timeOutId);
              for (const hex of msg) console.log(hex);

              console.log("");

              this.receiveResponse(msg);
              // this.client.close(() => {
              //   console.log(`${msg}\n CLOSED SOCKET`);
              // });
            });

            const closeSocketTimeout = setTimeout(() => {
              this.client.close(() => {
                console.log("Socket is closed");
              });
            }, 4999);

            const timeOutId = setTimeout(() => {
              this.client = dgram.createSocket("udp4");
              this.sendRequest(payload, requestType, packetNo, noOfPackets);
            }, this.timeout);
          }
        }
      }
    );
  }

  public sendResponse(
    id: string,
    msg: string,
    packetNo: number,
    noOfPackets: number
  ) {
    let header = constructHeaders(101, id, packetNo, noOfPackets);
    // Converts the message object into array
    const str = marshal(msg);
    const payload = Buffer.from(str + "\0");
    let bufferData = Buffer.concat([Buffer.from(header), payload]);
    this.client.send(
      bufferData,
      0,
      bufferData.length,
      this.sendPort,
      this.address,
      (err: Error | null, bytes: number) => {
        if (err) {
          console.log(`Error sending response: ${err}`);
        } else {
          console.log(`Sent ${bytes} to server`); //If successful Client will be notified...
          this.client.on("message", (message, rinfo) => {
            console.log("Response\n", message);
            this.pendingRequests.delete(id);

            // close the client socket
            // this.client.close(() => {
            //     console.log('Socket is closed');
            //   });
          });
        }
      }
    );
  }
}
