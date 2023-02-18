// | uint8, 1 byte: request type | string, 10 bytes: requestID | rest of payload

import {nanoid} from "nanoid";

export const craftHeaders = (requestType:number) => {
    let enc = new TextEncoder();
    const requestTypeByte = enc.encode(requestType.toString());
    const requestId = enc.encode(nanoid(10));

    const len = requestTypeByte.length + requestId.length;
    return Buffer.concat([requestTypeByte, requestId], len);

}

console.log(new TextDecoder().decode(craftHeaders(1)));