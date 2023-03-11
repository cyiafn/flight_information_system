export function convertToDateTime(time: bigint) {
  const timeInSeconds = (time * BigInt(1000)).toString();
  const d = new Date(Number(timeInSeconds));
  const dateFormatted = `${d.getDate()}-${
    d.getMonth() + 1
  }-${d.getFullYear()} ${d.getHours()}:${d.getMinutes()}`;

  return dateFormatted;
}

export function findStrFromBuffer(buffer: Buffer) {
  let idx = 0;
  for (const byte of buffer.values()) {
    if (byte == 0x00) break;
    idx++;
  }
  return { totalLen: idx + 1, str: buffer.toString("utf-8", 0, idx) };
}

export function getPacketInformation(packetId: string, packetNo: number, noOfPackets: number, requestType: number, payload: Buffer) {
  console.log("--------------------------------------------------------------")
  console.log(`The Packet Id is: ${packetId}.`)
  console.log(`This is Packet ${packetNo} out of ${noOfPackets}...`)
  console.log(`This is the Request Type ${requestType}`)
  console.log(`This is the Payload Received/Sent: ${payload}`)
  console.log("--------------------------------------------------------------")
}