// | uint8, 1 byte: request type | string, 10 bytes: requestID | uint8: Packet no. | uint8: no. of Packets | rest of payload

import {nanoid} from "nanoid";

export const createRequestId = () => {
    return nanoid(10);
}

export const craftHeaders = (requestType:number, requestIdStr : string, packetNo: number, noOfPackets: number) => {
    let enc = new TextEncoder();
    const requestTypeByte = enc.encode(requestType.toString());
    const requestId = enc.encode(requestIdStr);
    const requestPacketNo = enc.encode(packetNo.toString());
    const requestNoOfPackets = enc.encode(noOfPackets.toString());

    const len = requestTypeByte.length + requestId.length + requestPacketNo.length + requestNoOfPackets.length;
    
    return {id : requestIdStr , header: Buffer.concat([requestTypeByte, requestId, requestPacketNo,requestNoOfPackets], len)};

}