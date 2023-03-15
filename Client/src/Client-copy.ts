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
  promise: any;

  constructor(address: string, sendPort: number) {
    this.address = address; // IP Address of Server
    this.sendPort = sendPort; // Sending Port Number of Client
    this.receivePort = 0; // Client Port is using
    this.client = dgram.createSocket("udp4");
    this.requestId = ""; // Tracking of Request ID
    this.timeout = 5000; // Time out when no reply in 0.5 seconds
    this.monitorTimeOut = 0; // Time to set when to stop monitoring
    this.monitorMode = false; // Whether callback is enable
    this.promise = null;
  }

  private receiveResponse(buffer: Buffer) {
    const header = deconstructHeaders(buffer);

    logPacketInformation(
      header.requestId,
      Number(header.byteArrayBufferNo),
      Number(header.totalByteArrayBuffers),
      header.requestType,
      undefined
    );
    const payload = unmarshal(buffer.subarray(26, 512), header.requestType);

    if (typeof payload === "string") console.log(payload);
    else return payload;
  }

  public setMonitorTimeout(time: BigInt) {
    this.monitorTimeOut = Number(time);
  }
  //Send Method to cover both Idempotent and Non-Idempotent
  public async sendRequest(
    payload: Buffer,
    requestType: number,
    byteArrayBufferNo: number,
    totalByteArrayBuffers: number,
    responseLost = false
  ) {
    if (this.requestId === "") this.requestId = createRequestId();

    const header = constructHeaders(
      requestType,
      this.requestId,
      byteArrayBufferNo,
      totalByteArrayBuffers
    );

    // Converts the message object into array
    const buffer = Buffer.concat([header, payload], 512);

    this.promise = new Promise((resolve, reject) => {
      // Send over to server
      this.client.send(
        buffer,
        0,
        buffer.length,
        this.sendPort,
        this.address,
        (err: Error | null) => {
          if (err) {
            console.log(`Error sending message: ${err}`);
          } else {
            logPacketInformation(
              this.requestId,
              byteArrayBufferNo,
              totalByteArrayBuffers,
              requestType,
              payload
            );

            this.receivePort = this.client.address().port;
            this.client.on("message", (msg) => {
              let callback;

              // Simulate response Lost if true
              if (!responseLost) {
                clearTimeout(timeOutId);

                callback = this.receiveResponse(
                  msg
                ) as ResponseType.MonitorSeatUpdatesResponseType;
                responseLost = false;
              }

              // Dealing with callback
              if (callback === ResponseType.MonitorSeatUpdatesResponseType) {
                setTimeout(() => {
                  console.log("No more monitoring");
                  this.monitorMode = false;
                  this.client.close();
                  resolve(1);
                }, this.monitorTimeOut * 1000);
                this.monitorMode = true;
              } else if (!this.monitorMode) {
                this.client.close();
                resolve(1);
              }
            });

            // Resend if not acknowledgement has been received for 5 secs
            const timeOutId = setTimeout(() => {
              this.client = dgram.createSocket("udp4");
              console.log("No acknowledgement, sending packet again");
              this.sendRequest(
                payload,
                requestType,
                byteArrayBufferNo,
                totalByteArrayBuffers
              );
            }, this.timeout);
          }
        }
      );
    });
  }
}
