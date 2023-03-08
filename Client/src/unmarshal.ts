import { Buffer } from "buffer";
import { ResponseType, StatusCode } from "./interfaces";
import { convertToDateTime, findStrFromBuffer } from "./utility";

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
  let totalAvailableSeats: number, airfare: string, flightIdentifier: number;
  let departureTime: string;

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
      console.log(
        `These are the following Flight Identifier(s) from the given Source and Destination Location:`
      );
      for (const [idx, id] of flightIds.entries()) {
        if (idx === flightIds.length - 1) console.log(`${id}`);
        console.log(`${id}, `);
      }
      break;

    case ResponseType.GetFlightInformationResponseType:
      departureTime = convertToDateTime(buffer.readBigInt64LE());
      airfare = buffer.readDoubleLE(8).toFixed(2);
      totalAvailableSeats = buffer.readInt32LE(16);

      console.log(
        `The Departure for this Flight Identifier is ${departureTime}.`
      );
      console.log(`The Airfare cost ${airfare}.`);
      console.log(
        `The Total Available Seats for this Flight is ${totalAvailableSeats}.`
      );
      break;

    case ResponseType.MakeSeatReservationResponseType:
      // No Response Body
      console.log("Seats have been successfully reserved");
      break;

    case ResponseType.MonitorSeatUpdatesResponseType:
      return ResponseType.MonitorSeatUpdatesResponseType;

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

      departureTime = convertToDateTime(buffer.readBigInt64LE(curLen));

      airfare = buffer.readDoubleLE((curLen += 8)).toFixed(2);
      totalAvailableSeats = buffer.readInt32LE((curLen += 8));

      result = {
        FlightIdentifier: flightIdentifier,
        SourceLocation: sourceLocation,
        DestinationLocation: destinationLocation,
        DepartureTime: departureTime,
        Airfare: airfare,
        TotalAvailableSeats: totalAvailableSeats,
      };

      console.log(
        `The flight price for Flight Identifier ${flightIdentifier} flying from ${sourceLocation} to ${destinationLocation} at ${departureTime} have been updated.`
      );
      console.log(`The updated Airfare cost is ${airfare}.`);
      console.log(
        `The Total Available Seats for this Flight is ${totalAvailableSeats}.`
      );

      break;

    case ResponseType.CreateFlightResponseType:
      flightIdentifier = buffer.readInt32LE();
      console.log(
        `The Flight Identifier ${flightIdentifier} have been created.`
      );
      break;
  }
}

const buffer = Buffer.alloc(20);
buffer.writeBigInt64LE(BigInt(1678283107));
buffer.writeDoubleLE(1425.2, 8);
buffer.writeInt32LE(20, 16);

determineResponseType(buffer, 103);
