import dgram from 'dgram';
import { Buffer } from 'buffer';
import {
  constructHeaders,
  createRequestId,
  deconstructHeaders
} from './headers';
import { RequestObj, ResponseType } from './interfaces';
import { unmarshal } from './unmarshal';
import { clearTimeouts, logPacketInformation } from './utility';

const ac = new AbortController();

export class UDPClient {
  address: string;
  sendPort: number;
  receivePort: number;
  client: dgram.Socket;
  requestId: string;
  timeout: number;
  monitorTimeOut: number;
  retryCnt: number;
  maxRetries: number;
  monitorMode: boolean;
  timer: any;

  constructor(address: string, sendPort: number) {
    this.address = address; // IP Address of Server
    this.sendPort = sendPort; // Sending Port Number of Client
    this.receivePort = 0; // Client Port is using
    this.client = dgram.createSocket('udp4');
    this.requestId = ''; // Tracking of Request ID
    this.timeout = 5000; // Time out when no reply in 0.5 seconds
    this.monitorTimeOut = 0; // Time to set when to stop monitoring
    this.retryCnt = 0;
    this.maxRetries = 3;
    this.monitorMode = false;
    this.timer = [];
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

  private sendRequest(requestObj: RequestObj, buffer: Buffer) {
    return new Promise((resolve, reject) => {
      this.retryCnt++; // Increment the retry counter
      this.client.send(buffer, this.sendPort, this.address, (err) => {
        if (err) {
          console.log(`Error sending message: ${err}`);
          this.client.close();
          reject(err); // Reject the promise with the error
        } else {
          // Log information
          logPacketInformation(
            this.requestId,
            requestObj.byteArrayBufferNo,
            requestObj.totalByteArrayBuffers,
            requestObj.requestType,
            requestObj.payload
          );
          // Set a timer for 5 seconds
          this.timer.push(
            setTimeout(() => {
              console.log('No response received.');
              if (this.retryCnt < this.maxRetries) {
                console.log(
                  `Retrying... (attempt ${this.retryCnt + 1} of ${
                    this.maxRetries
                  })`
                );
                this.sendRequest(requestObj, buffer)
                  .then((response) => resolve(response)) // Resolve the promise with the response
                  .catch((err) => reject(err)); // Reject the promise with the error
              } else {
                console.log('else loop', this.retryCnt);
                console.log(
                  `Maximum number of retries (${this.maxRetries}) reached. Giving up. \n`
                );
                this.client.close();
                this.retryCnt = 0;
                reject(new Error('Maximum number of retries reached')); // Reject the promise with an error message
              }
            }, this.timeout)
          );
          if (
            (!requestObj.responseLost && this.retryCnt === 1) ||
            (requestObj.responseLost && this.retryCnt === 2)
          ) {
            this.client.on('message', (msg) => {
              const callback = this.receiveResponse(msg);
              clearTimeouts(this.timer);
              this.timer = [];
              if (callback === ResponseType.MonitorSeatUpdatesResponseType) {
                this.monitorMode = true;
                console.log(`Callback establised with server.`);
                const callbackTimer = setTimeout(() => {
                  console.log('Monitoring has expired...\n\n');
                  this.monitorMode = false;
                  resolve(1);
                  this.client.close();
                }, this.monitorTimeOut * 1000);
              } else if (!this.monitorMode) {
                resolve(1);
              }
            });
          }
        }
      });
    });
  }

  public sendRequests(requestObj: RequestObj) {
    const buffer = this.constructHeaderWithPayload(
      requestObj.payload,
      requestObj.requestType,
      requestObj.byteArrayBufferNo,
      requestObj.totalByteArrayBuffers
    );
    return new Promise(async (resolve, reject) => {
      try {
        await this.sendRequest(requestObj, buffer);
      } catch (e: any) {
        console.log(`Error: ${e.message} \n\n`);
      } finally {
        resolve(1);
      }
    });
  }
}
