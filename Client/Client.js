
const dgram = require('dgram');
const fakeStringify  = require('./fakeStringify.js');
const retrieveDataTypes = require('./retrieveDataTypes.js');
// Connect to Server Via UDP Connection
// port = 8080

class UDPClient {
    constructor(address, port)  // Specify the port number (agreed upon) and local ip address of the server.
    {
        this.address = address;
        this.port = port;
        this.client = dgram.createSocket('udp4');
        this.requestId = 0; //Tracking of Request ID
        this.pendingRequests = new Map();
        this.timeout = 10000; //Time out when no reply in 10 seconds  
        this.client.on('error', (err) => {
            this.handleError(); 
        })
        
        this.client.on('message', (message,rinfo) => { //client.on enables client to listen for message and execute handleMessage()
            this.handleMessage(msg,rinfo);
        })
    }


    //UDPClient get Methods
    getPort()
    {
        return this.port;
    }

    getAddress()
    {
        return this.address;
    }

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
        console.log('Client Error');
    }
    
    // TO-DO Method to create handle Request or Response...
    //handleMessage(msg,rinfo) 
   
    

    sendRequest(msg, isIdempotent = true) //Send Method to cover both Idempotent and Non-Idempotent
    {
        const id = this.requestId++;
        let stringData = fakeStringify({"id" : id, "type" : "request", "data" : msg}); // Converts the message object into array
        let variableArray = retrieveDataTypes(stringData); //Retrieve an string array of item types
        let bufferData = Buffer.from(variableArray + "," + stringData); //Converts the string into buffer data... 
        if (isIdempotent)
        {
            this.pendingRequests.set(id, {data, attempts: 0}) //Increase the pendingRequest if no response from server... Remove only when receive response...
        }

        this.client.send(bufferData,0, bufferData.length(), this.getPort(), this.getAddress(), (err, bytes) => {
            if (err) {
              console.log(`Error sending message: ${err}`);
            } else {
              console.log(`Sent ${bytes} bytes to server`); //If successful Client will be notified...
              if(!isIdempotent){
                setTimeout(()=> {
                    this.pendingRequests.delete(id); 
                }, this.getTimeout()); // Set timeout to the predefined timeout value
              }
            }
          });
    }

    sendResponse(id,msg) {
        let stringData = fakeStringify({"id" : id, "type" : "response", "data" : msg}); // Converts the message object into array
        let variableArray = retrieveDataTypes(stringData); //Retrieve an string array of item types
        let bufferData = Buffer.from(variableArray + "," + stringData); //Converts the string into buffer data... 
        this.client.send(bufferData,0, bufferData.length(), this.getPort(), this.getAddress(), (err, bytes) =>{
            if (err) {
                console.log(`Error sending response: ${err}`);
            }
            else{
                console.log(`Sent ${bytes} to server`); //If successful Client will be notified...
            }
        });
    }

}

//Testing if fake Stringify works...
// let result = (fakeStringify({ 'Id' : "1234", "type" : "request", 'Data' : [123,41,21,2] }));
// console.log(typeof(result))
// console.log(result);
// test = '{ "id" : "1234", "type" : "abc", "data" : {"id" : 55 } }'
// console.log(retrieveDataTypes(test));
let id = 123;
let msg = "hello";

let stringData = fakeStringify({"id" : id, "type" : "request", "data" : msg}); // Converts the message object into array
let variableArray = retrieveDataTypes(stringData);
let bufferData = Buffer.from(variableArray + "," + stringData);
console.log(stringData);
console.log(variableArray);
console.log(variableArray + "-" + stringData)
console.log(bufferData);
