export function convertToDateTime(time: bigint) {
  const timeInSeconds = (time * BigInt(1000)).toString();
  var d = new Date(Number(timeInSeconds));
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