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
      const lenOfArr = buffer.readBigInt64LE();
      const flightIds = [];
      let eSize = 8;
      for (let i = 0; i < lenOfArr; i++) {
        flightIds.push(buffer.readInt32LE(eSize));
        eSize += 4;
      }
      result = { FlightIdentifiers: flightIds };
      break;

    case ResponseType.GetFlightInformationResponseType:
      departureTime = buffer.readBigInt64LE();
      airfare = buffer.readDoubleLE(8);
      totalAvailableSeats = buffer.readInt32LE(16);
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
      totalAvailableSeats = buffer.readInt32LE();
      result = {
        TotalAvailableSeats: totalAvailableSeats,
      };
      break;

    case ResponseType.UpdateFlightPriceResponseType:
      flightIdentifier = buffer.readInt32LE();
      let curLen = 4;

      const { totalLen: totalLen1, str: sourceLocation } = findStrFromBuffer(
        buffer.subarray(4, buffer.length)
      );
      curLen += totalLen1;

      const { totalLen: totalLen2, str: destinationLocation } =
        findStrFromBuffer(buffer.subarray(curLen, buffer.length));
      curLen += totalLen2;

      departureTime = buffer.readBigInt64LE(curLen);

      airfare = buffer.readDoubleLE((curLen += 8));
      totalAvailableSeats = buffer.readInt32LE((curLen += 8));

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
      flightIdentifier = buffer.readInt32LE();
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
