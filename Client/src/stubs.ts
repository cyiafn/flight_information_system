import { UDPClient } from "./Client";
import { CreateFlightRequest, RequestType } from "./interfaces";
import { marshal } from "./marshal";

const client = new UDPClient("127.0.0.1", 8080);

// Get flight identifier
export async function getFlightIdentifier(
  sourceLocation: string,
  destinationLocation: string
) {
  const payload = marshal({
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
  });

  await client.sendRequest(
    payload,
    RequestType.GetFlightIdentifiersRequestType,
    1,
    1
  );
  console.log("HELLO");
}

// Get flight information
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

// Crate seat reservation
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

// Listen for seat updates
export function monitorSeatUpdatesCallbackRequest(
  flightIdentifier: number,
  lengthOfMonitorIntervalInSeconds: bigint
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    LengthOfMonitorIntervalInSeconds: lengthOfMonitorIntervalInSeconds,
  });

  client.sendRequest(payload, RequestType.MonitorSeatUpdatesRequestType, 1, 1);
  client.setMonitorTimeout(lengthOfMonitorIntervalInSeconds);
}

// Update Flight Price Request
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

// Create Flight Request
export function createFlightRequest(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1);
}

// Simulate Get Flight Information with request lost from server
export function createFlightWithRequestLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  console.log("Sending a packet, will resend in 5 sec...");
  setTimeout(
    () =>
      client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1),
    5000
  );
}

// Simulate Get Flight Information with response lost in client
export function createFlightWithResponseLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1, true);
}

// Simulate update Flight Price with request lost
export function updateFlightPriceRequestWithRequestLost(
  flightIdentifier: number,
  newPrice: number
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  console.log("Sending a packet, will resend in 5 sec...");
  setTimeout(
    () =>
      client.sendRequest(
        payload,
        RequestType.UpdateFlightPriceRequestType,
        1,
        1
      ),
    5000
  );
}

// Simulate update Flight Price with response lost
export function updateFlightPriceRequestWithResponseLost(
  flightIdentifier: number,
  newPrice: number
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  client.sendRequest(
    payload,
    RequestType.UpdateFlightPriceRequestType,
    1,
    1,
    true
  );
}
