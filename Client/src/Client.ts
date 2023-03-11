import dgram from "dgram";
import { Buffer } from "buffer";
import {
  constructHeaders,
  createRequestId,
  deconstructHeaders,
} from "./headers";
import { ResponseType } from "./interfaces";
import { unmarshal } from "./unmarshal";
import { logPacketInformation } from "./utility";

export class UDPClient {
  address: string;
  sendPort: number;
  receivePort: number;
  client: dgram.Socket;
  requestId: string;
  timeout: number;
  monitorTimeOut: number;
  monitorMode: boolean;

  constructor(address: string, sendPort: number) {
    this.address = address; // IP Address of Server
    this.sendPort = sendPort; // Sending Port Number of Client
    this.receivePort = 0; // Client Port is using
    this.client = dgram.createSocket("udp4");
    this.requestId = ""; // Tracking of Request ID
    this.timeout = 5000; // Time out when no reply in 0.5 seconds
    this.monitorTimeOut = 0; // Time to set when to stop monitoring
    this.monitorMode = false; // Whether callback is enable
  }

  private receiveResponse(buffer: Buffer) {
    const header = deconstructHeaders(buffer);

    logPacketInformation(
      header.requestId,
      Number(header.packetNo),
      Number(header.noOfPackets),
      header.requestType,
      undefined
    );
    const payload = unmarshal(buffer.subarray(26, 512), header.requestType);

    if (typeof payload === "string") console.log(payload);
    else return payload;
  }

  //Send Method to cover both Idempotent and Non-Idempotent
  public sendRequest(
    payload: Buffer,
    requestType: number,
    packetNo: number,
    noOfPackets: number
  ) {
    this.requestId = createRequestId();
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
          logPacketInformation(
            this.requestId,
            packetNo,
            noOfPackets,
            requestType,
            payload
          );

          this.receivePort = this.client.address().port;
          this.client.on("message", (msg) => {
            clearTimeout(closeSocketTimeout);
            clearTimeout(timeOutId);

            const callback = this.receiveResponse(
              msg
            ) as ResponseType.MonitorSeatUpdatesResponseType;

            // Dealing with callback
            if (callback === ResponseType.MonitorSeatUpdatesResponseType) {
              setTimeout(() => {
                console.log("No more monitoring");
                this.monitorMode = false;
                this.client.close();
              }, this.monitorTimeOut * 1000);
              this.monitorMode = true;
            } else if (!this.monitorMode)
              this.client.close(() => {
                console.log(`Socket is closed after receiving acknowledgement`);
              });
          });

          const closeSocketTimeout = setTimeout(() => {
            this.client.close(() => {
              console.log(
                "Socket is closed after trying to send same packet again"
              );
            });
          }, this.timeout - 1);

          const timeOutId = setTimeout(() => {
            this.client = dgram.createSocket("udp4");
            this.sendRequest(payload, requestType, packetNo, noOfPackets);
          }, this.timeout);
        }
      }
    );
  }
}
