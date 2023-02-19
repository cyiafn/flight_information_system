
import dgram from 'dgram';
import { Buffer } from 'buffer';
import { marshall } from './marshalling';
import { retrieveDataTypes } from './retrieveDataTypes';
import { craftHeaders } from './headers';
import { buffer } from 'stream/consumers';
import { RemoteInfo } from 'dgram';

// Connect to Server Via UDP Connection
// port = 8080


class UDPClient {
    address: string;
    sendPort: number;
    receivePort: number;
    client: dgram.Socket;
    requestId: number;
    pendingRequests: Map<any, any>;
    timeout: number;
    constructor(address:string, sendPort:number, receivePort : number)  // Specify the port number (agreed upon) and local ip address of the server.
    {
        this.address = address; //IP Address of Server
        this.sendPort = sendPort; //Sending Port Number of Client
        this.receivePort = receivePort; //Listening Port of Client
        this.client = dgram.createSocket('udp4');
        this.requestId = 0; //Tracking of Request ID
        this.pendingRequests = new Map(); //Storing of Pending Requests
        this.timeout = 5000; //Time out when no reply in 0.5 seconds  

        this.client.bind(this.receivePort, this.address);
        
        this.client.on('listening', () => {
            console.log(`test Client listening on ${this.address}:${this.receivePort}`);
        });

    }
    // handleMessage(msg: string, rinfo: dgram.RemoteInfo) {
    //     throw new Error('Method not implemented.');
    // }


    //UDPClient get Methods

    getRequestId()
    {
        return this.requestId;
    }
    
    getPendingRequests()
    {
        return this.pendingRequests;
    }
    
    getTimeout()
    {
        return this.timeout;
    }

    handleError() //Method to handle error -> TO-DO Handle Specific Error by their error code
    {
        console.log('Client Encountered an Error');
    }
    
    // TO-DO Method to create handle Request or Response...
    handleMessage(msg:string, rinfo: RemoteInfo) {
        console.log("ACK");
    }


    

    //Send Method to cover both Idempotent and Non-Idempotent
    public sendRequest(msg:string) 
    {
        let id = this.requestId;
        let header = craftHeaders(this.requestId++);
        // Converts the message object into array
        const {str, attr} = marshall(msg); 
        // 0 = "string" | 1 = "number" | 2 = "boolean" | 3 = "array" | 4 = "object"
        
        // 6[4,0,0,0,0,0] | payload

        //Converts the string into buffer data with concat of headers
        const payload = Buffer.from(attr.toString() + "|" + str + "\0");
        let bufferData = Buffer.concat([Buffer.from(header), payload]);
        this.pendingRequests.set(id, {bufferData, attempts: 0})

        

        this.client.send(bufferData, 0, bufferData.length, this.sendPort, this.address, (err: Error | null, bytes : number) => {
            if (err) {
                console.log(`Error sending message: ${err}`);
            } else {
                //If successful Client will be notified...
                console.log(`Sent ${bytes} bytes to server`); 

                const timeOutId = setTimeout(()=> {
                    this.sendRequest(msg);
                }, this.getTimeout());

                this.client.on('message', (message, rinfo) => {
                    console.log("ack" + message);
                    clearTimeout(timeOutId);
                    this.pendingRequests.delete(id);
                    
                    
                    // check if there are any more pending requests before closing the socket..
                    if (this.pendingRequests.size === 0) {
                      // close the socket
                      
                      this.client.close(() => {
                        console.log('Socket is closed');
                      });
                    }
                  });
            }
        });

        // Send an empty packet once a request is completed...
        const emptyBuffer = Buffer.alloc(0);
        this.client.send(emptyBuffer, 0, 0, this.sendPort, this.address, (err) => {
        if (err) {
            console.error('Error sending empty packet:', err);
        } else {
            console.log('Empty packet sent successfully');
        }
        });
    }

    public sendResponse(id:string,msg:string) {
        let header = craftHeaders(101);
        // Converts the message object into array
        const {str, attr} = marshall(msg); 
        // 0 = "string" | 1 = "number" | 2 = "boolean" | 3 = "array" | 4 = "object"
        const payload = Buffer.from(attr.toString() + "|" + str + "\0");
        let bufferData = Buffer.concat([Buffer.from(header), payload]);
        this.client.send(bufferData,0, bufferData.length, this.sendPort, this.address, (err : Error | null, bytes : number) =>{
            if (err) {
                console.log(`Error sending response: ${err}`);
            }
            else{
                console.log(`Sent ${bytes} to server`); //If successful Client will be notified...
                this.client.on('message', (message, rinfo) => {
                    console.log("ack" + message);
                    this.pendingRequests.delete(id);
                    
                    // close the client socket
                    this.client.close(() => {
                        console.log('Socket is closed');
                      });
                  });  
            }
        })
    }

}

let id = 123;
let msg = "hello";
let msg2 = "Can i book this?"
// let header = craftHeaders(100);
// let data = {data : [1,2,3,4,5]}
// const {str, attr} = marshall(msg);  // Converts the message object into array
// let dataString = attr.toString() + "|" + str + "\0";
// console.log(dataString);

let client = new UDPClient('127.0.0.1', 3333, 4444);
client.sendRequest(msg);
// client.sendRequest(msg2)
// client.sendRequest("Booking Flight: 1234");
