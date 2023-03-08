import { Buffer } from "buffer";
import { ResponseType, StatusCode } from "./interfaces";

export function unmarshal(buffer: Buffer, requestType: Number) {
  const statusCode = buffer[0];
  const data = buffer.subarray(1, buffer.length);

  switch (statusCode) {
    case StatusCode.BusinessLogicGenericError:
      return "Generic Error";
    case StatusCode.MarshallerError:
      return "Marshal Error";
    case StatusCode.NoMatchForSourceAndDestination:
      return "No match for source and destination";
    case StatusCode.NoSuchFlightIdentifier:
      return "No flight identifier";
    case StatusCode.InsufficientNumberOfAvailableSeats:
      return "Insufficient number of available seats";
    case StatusCode.Success:
      return determineResponseType(data, requestType);
  }
}

function determineResponseType(buffer: Buffer, responseType: Number) {
  let result;
  let totalAvailableSeats: number, airfare: number, flightIdentifier: number;
  let departureTime: BigInt;

  switch (responseType) {
    case ResponseType.PingResponseType:
      console.log("PONG");
      break;

    case ResponseType.GetFlightIdentifiersResponseType:
      const lenOfArr = buffer.readBigUInt64LE();
      const flightIds = [];
      let eSize = 8;
      for (let i = 0; i < lenOfArr; i++) {
        flightIds.push(buffer.readInt32LE(eSize));
        eSize += 4;
      }
      result = { FlightIdentifiers: flightIds };
      break;

    case ResponseType.GetFlightInformationResponseType:
      departureTime = buffer.readBigUInt64LE();
      airfare = buffer.readDoubleLE(8);
      totalAvailableSeats = buffer.readUInt32LE(16);
      result = {
        DepartureTime: departureTime,
        Airfare: airfare,
        TotalAvailableSeats: totalAvailableSeats,
      };
      break;

    case ResponseType.MakeSeatReservationResponseType:
      // No Response Body
      break;

    case ResponseType.MonitorSeatUpdatesResponseType:
      return "Monitor Success";
      break;

    case ResponseType.MonitorSeatUpdatesCallbackType:
      totalAvailableSeats = buffer.readUInt32LE();
      result = {
        TotalAvailableSeats: totalAvailableSeats,
      };
      break;

    case ResponseType.UpdateFlightPriceResponseType:
      flightIdentifier = buffer.readUInt32LE();
      let curLen = 4;

      const { totalLen: totalLen1, str: sourceLocation } = findStrFromBuffer(
        buffer.subarray(4, buffer.length)
      );
      curLen += totalLen1;

      const { totalLen: totalLen2, str: destinationLocation } =
        findStrFromBuffer(buffer.subarray(curLen, buffer.length));
      curLen += totalLen2;

      departureTime = buffer.readBigUInt64LE(curLen);

      airfare = buffer.readDoubleLE((curLen += 8));
      totalAvailableSeats = buffer.readUint32LE((curLen += 8));

      result = {
        FlightIdentifier: flightIdentifier,
        SourceLocation: sourceLocation,
        DestinationLocation: destinationLocation,
        DepartureTime: departureTime,
        Airfare: airfare,
        TotalAvailableSeats: totalAvailableSeats,
      };
      break;

    case ResponseType.CreateFlightResponseType:
      flightIdentifier = buffer.readUInt32LE();
      return { FlightIdentifier: flightIdentifier };
  }

  return result;
}

function findStrFromBuffer(buffer: Buffer) {
  let idx = 0;
  for (const byte of buffer.values()) {
    if (byte == 0x00) break;
    idx++;
  }
  return { totalLen: idx + 1, str: buffer.toString("utf-8", 0, idx) };
}

//Case 102
// const buffer = Buffer.alloc(20);
// buffer.writeBigUInt64LE(BigInt(3));
// buffer.writeUint32LE(123, 8);
// buffer.writeUint32LE(45, 12);
// buffer.writeUint32LE(100, 16);
// console.log(buffer);
// const ans = determineResponseType(buffer, 102);
// console.log(ans);

//Case 103
// const buffer = Buffer.alloc(20);
// buffer.writeBigUInt64LE(BigInt(1400), 0);
// buffer.writeFloatLE(12000.25, 8);
// buffer.writeUint32LE(100, 16);
// console.log(buffer);
// const ans = determineResponseType(buffer, 103);
// console.log(ans);

//Case 105
// const buffer = Buffer.alloc(4);
// buffer.writeUint32LE(100, 0);
// console.log(buffer);
// const ans = determineResponseType(buffer, 105);
// console.log(ans);

//Case 106
// const buffer = Buffer.alloc(43);
// buffer.writeUint32LE(100);
// buffer.write("Singapore\0", 4);
// buffer.write("Malaysia\0", 14);
// buffer.writeBigUint64LE(BigInt(1400), 23);
// buffer.writeFloatLE(140.99, 31);
// buffer.writeUint32LE(69, 39);
// const ans = determineResponseType(buffer, 106);
// console.log(ans);

//case 107
// const buffer = Buffer.alloc(4);
// buffer.writeUint32LE(150);
// const ans = determineResponseType(buffer, 107);
// console.log(ans);
