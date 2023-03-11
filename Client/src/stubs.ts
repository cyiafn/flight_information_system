import { UDPClient } from "./Client";
import { CreateFlightRequest, RequestType } from "./interfaces";
import { marshal } from "./marshal";

const client = new UDPClient("127.0.0.1", 8080);

export function getFlightIdentifier(
  sourceLocation: string,
  destinationLocation: string
) {
  const payload = marshal({
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
  });

  client.sendRequest(
    payload,
    RequestType.GetFlightIdentifiersRequestType,
    1,
    1
  );
}

export function getFlightInformation(flightIdentifier: number) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
  });

  client.sendRequest(
    payload,
    RequestType.GetFlightInformationRequestType,
    1,
    1
  );
}

export function createSeatReservationRequest(
  flightIdentifier: number,
  seatsToReserve: number
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    SeatsToReserve: seatsToReserve,
  });

  client.sendRequest(payload, RequestType.MakeSeatReservationRequestType, 1, 1);
}

export function monitorSeatUpdatesCallbackRequest(
  flightIdentifier: number,
  lengthOfMonitorIntervalInSeconds: bigint
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    LengthOfMonitorIntervalInSeconds: lengthOfMonitorIntervalInSeconds,
  });

  client.sendRequest(payload, RequestType.MonitorSeatUpdatesRequestType, 1, 1);
}

export function updateFlightPriceRequest(
  flightIdentifier: number,
  newPrice: number
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  client.sendRequest(payload, RequestType.UpdateFlightPriceRequestType, 1, 1);
}

export function createFlightRequest(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1);
}

