import dgram from 'dgram';
import { Buffer } from 'buffer';
import {
  constructHeaders,
  createRequestId,
  deconstructHeaders
} from './headers';
import { RequestObj, ResponseType } from './interfaces';
import { unmarshal } from './unmarshal';
import { isTimeout, logPacketInformation } from './utility';

const ac = new AbortController();

export class UDPClient {
  address: string;
  sendPort: number;
  receivePort: number;
  client: dgram.Socket;
  requestId: string;
  timeout: number;
  monitorTimeOut: number;
  maxRequests: number;
  reply: boolean;
  timer: any;

  constructor(address: string, sendPort: number) {
    this.address = address; // IP Address of Server
    this.sendPort = sendPort; // Sending Port Number of Client
    this.receivePort = 0; // Client Port is using
    this.client = dgram.createSocket('udp4');
    this.requestId = ''; // Tracking of Request ID
    this.timeout = 5000; // Time out when no reply in 0.5 seconds
    this.monitorTimeOut = 0; // Time to set when to stop monitoring
    this.maxRequests = 4;
    this.reply = false;
    this.timer = null;
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

    if (typeof payload === 'string') console.log(payload);
    else return payload;
  }

  private constructHeaderWithPayload(
    payload: Buffer,
    requestType: number,
    byteArrayBufferNo: number,
    totalByteArrayBuffers: number
  ) {
    if (this.requestId === '') this.requestId = createRequestId();

    const header = constructHeaders(
      requestType,
      this.requestId,
      byteArrayBufferNo,
      totalByteArrayBuffers
    );

    // Converts the message object into array
    return Buffer.concat([header, payload], 512);
  }

  private sendRequest(buffer: Buffer) {
    return new Promise((resolve, reject) => {
      this.client.send(
        buffer,
        0,
        buffer.length,
        this.sendPort,
        this.address,
        (err: Error | null) => {
          if (err) {
            reject(`Error sending message: ${err}`);
          } else {
            resolve('Sent!');
          }
        }
      );
    });
  }

  public async sendRequests(requestObj: RequestObj) {
    this.client.on('message', (msg) => {
      this.receiveResponse(msg);
      this.reply = true;
    });

    const buffer = this.constructHeaderWithPayload(
      requestObj.payload,
      requestObj.requestType,
      requestObj.byteArrayBufferNo,
      requestObj.totalByteArrayBuffers
    );
    return new Promise(async (resolve, reject) => {
      for (let i = 0; i < this.maxRequests; i++) {
        if (i > 0)
          console.log(
            `No acknowledgement from server, sending ${i} / 3 retries`
          );
        logPacketInformation(
          this.requestId,
          requestObj.byteArrayBufferNo,
          requestObj.totalByteArrayBuffers,
          requestObj.requestType,
          requestObj.payload
        );

        await this.sendRequest(buffer);
        await isTimeout(this.timeout);
        if (this.reply) {
          break;
        }
      }

      this.client.close();
      if (this.reply) {
        resolve(1);
        this.reply = false;
      } else reject('Exceeded the requests limit...\n');
    });
  }

  public async sendRequestCallback(requestObj: RequestObj) {
    return new Promise(async (resolve, reject) => {
      this.client.on('message', (msg) => {
        const callback = this.receiveResponse(msg);
        if (callback === ResponseType.MonitorSeatUpdatesResponseType) {
          this.reply = true;
          setTimeout(() => {
            console.log('No more monitoring');
            this.client.close();
            resolve(1);
            this.reply = false;
          }, this.monitorTimeOut * 1000);
        }
      });

      const buffer = this.constructHeaderWithPayload(
        requestObj.payload,
        requestObj.requestType,
        requestObj.byteArrayBufferNo,
        requestObj.totalByteArrayBuffers
      );
      for (let i = 0; i < this.maxRequests; i++) {
        if (i > 0)
          console.log(
            `No acknowledgement from server, sending ${i} / 3 retries`
          );
        logPacketInformation(
          this.requestId,
          requestObj.byteArrayBufferNo,
          requestObj.totalByteArrayBuffers,
          requestObj.requestType,
          requestObj.payload
        );

        await this.sendRequest(buffer);
        await isTimeout(this.timeout);
        if (this.reply) {
          console.log('Monitoring now...');
          break;
        }
      }
      if (!this.reply) reject('Exceeded the requests limit...\n');
    });
  }
}
