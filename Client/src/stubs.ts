import { UDPClient } from "./Client";
import { CreateFlightRequest, RequestType } from "./interfaces";
import { marshal } from "./marshal";

const ip = process.env.IP || "localhost";

// Get flight identifier
export async function getFlightIdentifier(
  sourceLocation: string,
  destinationLocation: string
) {
  const payload = marshal({
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(
      payload,
      RequestType.GetFlightIdentifiersRequestType,
      1,
      1
    );
    resolve(client.promise);
  });
}

// Get flight information
export function getFlightInformation(flightIdentifier: number) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(
      payload,
      RequestType.GetFlightInformationRequestType,
      1,
      1
    );
    resolve(client.promise);
  });
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

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(
      payload,
      RequestType.MakeSeatReservationRequestType,
      1,
      1
    );
    resolve(client.promise);
  });
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

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(
      payload,
      RequestType.MonitorSeatUpdatesRequestType,
      1,
      1
    );
    client.setMonitorTimeout(lengthOfMonitorIntervalInSeconds);
    resolve(client.promise);
  });
}

// Update Flight Price Request
export function updateFlightPriceRequest(
  flightIdentifier: number,
  newPrice: string
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(payload, RequestType.UpdateFlightPriceRequestType, 1, 1);
    resolve(client.promise);
  });
}

// Create Flight Request
export function createFlightRequest(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1);
    resolve(client.promise);
  });
}

// Simulate Get Flight Information with request lost from server
export function createFlightWithRequestLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    console.log("Sending a packet, will resend in 5 sec...");
    setTimeout(() => {
      client.sendRequest(payload, RequestType.CreateFlightRequestType, 1, 1);
      resolve(client.promise);
    }, 5000);
  });
}

// Simulate Get Flight Information with response lost in client
export function createFlightWithResponseLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.sendRequest(
      payload,
      RequestType.CreateFlightRequestType,
      1,
      1,
      true
    );
    resolve(client.promise);
  });
}

// Simulate update Flight Price with request lost
export function updateFlightPriceRequestWithRequestLost(
  flightIdentifier: number,
  newPrice: string
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    console.log("Sending a packet, will resend in 5 sec...");
    setTimeout(() => {
      client.sendRequest(
        payload,
        RequestType.UpdateFlightPriceRequestType,
        1,
        1
      );
      resolve(client.promise);
    }, 5000);
  });
}

// Simulate update Flight Price with response lost
export function updateFlightPriceRequestWithResponseLost(
  flightIdentifier: number,
  newPrice: string
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice,
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    setTimeout(() => {
      client.sendRequest(
        payload,
        RequestType.UpdateFlightPriceRequestType,
        1,
        1
      );
      resolve(client.promise);
    }, 5000);
  });
}
