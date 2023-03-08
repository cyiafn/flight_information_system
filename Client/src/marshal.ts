import { Buffer } from "buffer";
import {
  instanceOfCreateFlightRequest,
  instanceOfGetFlightIdentifiersRequest,
  instanceOfGetFlightInformationRequest,
  instanceOfMakeSeatReservationRequest,
  instanceOfMonitorSeatUpdatesRequest,
  instanceOfUpdateFlightPriceRequest,
} from "./interfaces";

export function marshal(data: any): any {
  let requestBuffer;
  if (instanceOfGetFlightIdentifiersRequest(data)) {
    const lengthReq =
      data.SourceLocation.length + data.DestinationLocation.length + 2;
    requestBuffer = Buffer.alloc(lengthReq);
    requestBuffer.write(data.SourceLocation);
    requestBuffer.write("\0", data.SourceLocation.length);
    requestBuffer.write(
      data.DestinationLocation,
      data.SourceLocation.length + 1
    );
    requestBuffer.write("\0", lengthReq - 1);
  } else if (instanceOfGetFlightInformationRequest(data)) {
    requestBuffer = Buffer.alloc(4);
    requestBuffer.writeUint32LE(data.FlightIdentifier);
  } else if (instanceOfMakeSeatReservationRequest(data)) {
    requestBuffer = Buffer.alloc(8);
    requestBuffer.writeUint32LE(data.FlightIdentifier);
    requestBuffer.writeUint32LE(data.SeatsToReserve, 4);
  } else if (instanceOfMonitorSeatUpdatesRequest(data)) {
    requestBuffer = Buffer.alloc(12);
    requestBuffer.writeUint32LE(data.FlightIdentifier);
    requestBuffer.writeBigUInt64LE(
      BigInt(data.LengthOfMonitorIntervalInSeconds),
      4
    );
  } else if (instanceOfUpdateFlightPriceRequest(data)) {
    requestBuffer = Buffer.alloc(12);
    requestBuffer.writeUint32LE(data.FlightIdentifier);
    requestBuffer.writeDoubleLE(data.NewPrice, 4);
  } else if (instanceOfCreateFlightRequest(data)) {
    const lengthReq =
      data.SourceLocation.length +
      data.DestinationLocation.length +
      2 +
      8 +
      8 +
      4;
    let start = data.SourceLocation.length;
    requestBuffer = Buffer.alloc(lengthReq);
    requestBuffer.write(data.SourceLocation);
    requestBuffer.write("\0", start);

    requestBuffer.write(data.DestinationLocation, ++start);
    requestBuffer.write("\0", (start += data.DestinationLocation.length));

    requestBuffer.writeBigUInt64LE(BigInt(data.DepartureTime), ++start);

    requestBuffer.writeDoubleLE(data.Airfare, (start += 8));

    requestBuffer.writeUInt32LE(data.TotalAvailableSeats, (start += 8));
  }

  return requestBuffer;
}
