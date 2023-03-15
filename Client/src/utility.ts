// Utility function to convert date-time
export function convertToDateTime(time: bigint) {
  const timeInSeconds = (time * BigInt(1000)).toString();
  const d = new Date(Number(timeInSeconds));
  const dateFormatted = `${d.getDate()}-${
    d.getMonth() + 1
  }-${d.getFullYear()} ${d.getUTCHours()}H:${d.getMinutes()}M`;

  return dateFormatted;
}

// Utility function to find str datatype from buffer
export function findStrFromBuffer(buffer: Buffer) {
  let idx = 0;
  for (const byte of buffer.values()) {
    if (byte == 0x00) break;
    idx++;
  }
  return { totalLen: idx + 1, str: buffer.toString('utf-8', 0, idx) };
}

// Print out packet information
export function logPacketInformation(
  packetId: string,
  packetNo: number,
  noOfPackets: number,
  requestType: number,
  payload: Buffer | undefined
) {
  console.log('--------------------Packet Information-----------------------');
  console.log(`The Packet Id is: ${packetId}`);
  console.log(`This is Packet ${packetNo} out of ${noOfPackets}`);

  if (requestType >= 2 && requestType <= 7) {
    console.log(`This is the Request Type ${requestType}`);
    console.log(`This is the Payload Sent: ${payload?.toString('hex')}`);
  } else console.log(`This is the Response Type ${requestType}`);
  console.log(
    '--------------------------------------------------------------\n'
  );
}

export function clearTimeouts(timeoutArr: any) {
  for (var i = 0; i < timeoutArr.length; i++) {
    clearTimeout(timeoutArr[i]);
  }
}

export function isTimeout(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time));
}
