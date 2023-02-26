
import dgram, { RemoteInfo } from 'dgram';
import { Buffer } from 'buffer';
import { marshal } from './marshal';
import { craftHeaders, createRequestId } from './headers';
import { PendingRequest } from './interfaces';

// Connect to Server Via UDP Connection
// port = 8080


class UDPClient {
    address: string;
    sendPort: number;
    receivePort: number;
    client: dgram.Socket;
    requestId: string;
    pendingRequests: Map<string, PendingRequest>;
    timeout: number;
    constructor(address:string, sendPort:number, receivePort : number)  // Specify the port number (agreed upon) and local ip address of the server.
    {
        this.address = address; //IP Address of Server
        this.sendPort = sendPort; //Sending Port Number of Client
        this.receivePort = receivePort; //Listening Port of Client
        this.client = dgram.createSocket('udp4');
        this.requestId = ''; //Tracking of Request ID
        this.pendingRequests = new Map(); //Storing of Pending Requests
        this.timeout = 5000; //Time out when no reply in 0.5 seconds  

        this.client.bind(this.receivePort, this.address);
        
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
    getRequestId() {
        return this.requestId;
    }
    
    setPendingRequests(id:string , bufferData: string) {
        let attempt = 0;
        if(this.pendingRequests.has(id))
            attempt = (this.pendingRequests.get(id)?.attempts || 0 ) + 1;
            
        this.pendingRequests.set(id, {data: bufferData, attempts: attempt})
    }

    getPendingRequests(id:string) {
        return this.pendingRequests.get(id);
    }
    
    getTimeout() {
        return this.timeout;
    }

    //Method to handle error -> TO-DO Handle Specific Error by their error code
    handleError() {
        console.log('Client Encountered an Error');
    }
    
    // TO-DO Method to create handle Request or Response...
    handleMessage(msg:string, rinfo: RemoteInfo) {
        console.log("ACK");
    }
    
    public sendMultipleRequests(msg:string)
    {
        // 13 bytes header message max 498 Bytes  1 byte terminate
        let bytesStream = new TextEncoder().encode(msg);
        const decoder = new TextDecoder()
        let noOfPackets = Math.ceil(bytesStream.length / 498);
        let packetNo = 1;

        // Create Request Id for these spliced msg requests
        const requestIdStr = createRequestId();

        while (bytesStream.length) {
            const maxBytes = Math.min(bytesStream.length, 498);
            const slicedMsg = decoder.decode(bytesStream.slice(0, maxBytes));
            
            this.sendRequest(slicedMsg, requestIdStr, packetNo++, noOfPackets);

            bytesStream = bytesStream.slice(maxBytes);
        }
    }

    //Send Method to cover both Idempotent and Non-Idempotent
    public sendRequest(msg:string, requestIdStr: string, packetNo: number, noOfPackets: number) 
    { 
        const {id, header} = craftHeaders(1,requestIdStr, packetNo, noOfPackets);
        this.setRequestId(id);
        // Converts the message object into array
        const str = marshal(msg); 

        //Converts the string into buffer data with concat of headers
        const payload = Buffer.from(str + "\0");
        let bufferData = Buffer.concat([Buffer.from(header), payload]);
        
        this.client.send(bufferData, 0, bufferData.length, this.sendPort, this.address, (err: Error | null, bytes : number) => {
            if (err) {
                console.log(`Error sending message: ${err}`);
            } else {
                //If successful Client will be notified...
                console.log(`Sent ${bytes} bytes to server`); 

                const timeOutId = setTimeout(()=> {
                    this.sendRequest(msg, requestIdStr , packetNo, noOfPackets);
                }, this.getTimeout());

                this.client.on('message', (message, rinfo) => {
                    console.log(rinfo);
                    clearTimeout(timeOutId);
                    this.pendingRequests.delete(id);

                    
                    // check if there are any more pending requests before closing the socket..
                    // if (this.pendingRequests.size === 0) {
                    //   this.client.close(() => {
                    //     console.log('Socket is closed');
                    //   });
                    // }
                  });
            }
        });

        // Send an empty packet once a request is completed...
        const emptyBuffer = Buffer.alloc(512);
        this.client.send(emptyBuffer, 0, emptyBuffer.length, this.sendPort, this.address, (err) => {
        if (err) {
            console.error('Error sending empty packet:', err);
        } else {
            console.log(`Empty packet ${emptyBuffer.byteLength} bytes sent successfully`);
        }
        });
    }

    public sendResponse(id:string, msg:string, packetNo: number, noOfPackets: number) {
        let {id: _, header} = craftHeaders(101, id, packetNo, noOfPackets);
        // Converts the message object into array
        const str = marshal(msg); 
        const payload = Buffer.from(str + "\0");
        let bufferData = Buffer.concat([Buffer.from(header), payload]);
        this.client.send(bufferData,0, bufferData.length, this.sendPort, this.address, (err : Error | null, bytes : number) =>{
            if (err) {
                console.log(`Error sending response: ${err}`);
            }
            else{
                console.log(`Sent ${bytes} to server`); //If successful Client will be notified...
                this.client.on('message', (message, rinfo) => {
                    console.log("Response\n", message);
                    this.pendingRequests.delete(id);
                    
                    // close the client socket
                    // this.client.close(() => {
                    //     console.log('Socket is closed');
                    //   });
                  });  
            }
        })
    }

}

let id = 123;
let msg = "hello";
// let msg2 = marshal({"Booking":1234, "seat":23, "Name" : "Eric"})
// let msg3 = marshal({"Booking":1234, "seat":[23,24,25], "Name" : "Mr."})


let client = new UDPClient('127.0.0.1', 3333, 4444);
const fakeMsg = 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab'
client.sendMultipleRequests(fakeMsg);
