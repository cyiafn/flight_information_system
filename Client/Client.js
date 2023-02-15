
const dgram = require('dgram');
const Stringify  = require('./fakeStringify.js');

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
        
        this.client.on('message', (message,rinfo) => {
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
    
    //handleMessage(msg,rinfo) // TO-DO Method to create handle Request or Response...
   
    

    send(msg, isIdempotent) //Send Method to cover both Idempotent and Non-Idempotent
    {
        const id = this.requestId++;
        let data = Buffer.from(Stringify(msg)); // Converts the message into buffer data... TO-DO Create a Stringify Function to convert to {"code" : xxx , "data" : {}}
        
        if (isIdempotent)
        {
            this.pendingRequests.set(id, {data, attempts: 0})
        }

        this.client.send(data, data.length, this.getPort(), this.getAddress(), (err, bytes) => {
            if (err) {
              console.log(`Error sending response: ${err}`);
            } else {
              console.log(`Sent ${bytes} bytes to server`);
            }
          });
    }

}

//Testing if fake Stringify works...
console.log(Stringify({"Data" : [123,41,21,2]}))