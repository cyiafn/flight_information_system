import { UDPClient } from './Client';
import { CreateFlightRequest, RequestType } from './interfaces';
import { marshal } from './marshal';

const ip = process.env.IP || 'localhost';

// Get flight identifier
export async function getFlightIdentifier(
  sourceLocation: string,
  destinationLocation: string
) {
  const payload = marshal({
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.GetFlightIdentifiersRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1
    });
    resolve(1);
  });
}

// Get flight information
export function getFlightInformation(flightIdentifier: number) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    resolve(
      client.sendRequests({
        payload: payload,
        requestType: RequestType.GetFlightInformationRequestType,
        byteArrayBufferNo: 1,
        totalByteArrayBuffers: 1
      })
    );
  });
}

// Crate seat reservation
export function createSeatReservationRequest(
  flightIdentifier: number,
  seatsToReserve: number
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    SeatsToReserve: seatsToReserve
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.MakeSeatReservationRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1
    });
    resolve(1);
  });
}

// Listen for seat updates
export function monitorSeatUpdatesCallbackRequest(
  flightIdentifier: number,
  lengthOfMonitorIntervalInSeconds: bigint
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    LengthOfMonitorIntervalInSeconds: lengthOfMonitorIntervalInSeconds
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    client.monitorTimeOut = Number(lengthOfMonitorIntervalInSeconds);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.MonitorSeatUpdatesRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1
    });
    resolve(1);
  });
}

// Update Flight Price Request
export function updateFlightPriceRequest(
  flightIdentifier: number,
  newPrice: string
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.UpdateFlightPriceRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1
    });
    resolve(1);
  });
}

// Create Flight Request
export function createFlightRequest(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.CreateFlightRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1
    });
    resolve(1);
  });
}

// Simulate Get Flight Information with request lost from server
export function createFlightWithRequestLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    console.log('Sending a packet, will resend in 5 sec...');
    setTimeout(async () => {
      await client.sendRequests({
        payload: payload,
        requestType: RequestType.CreateFlightRequestType,
        byteArrayBufferNo: 1,
        totalByteArrayBuffers: 1
      });
    }, 5000);
  });
}

// Simulate Get Flight Information with response lost in client
export function createFlightWithResponseLost(dto: CreateFlightRequest) {
  const payload = marshal({
    ...dto
  });

  return new Promise(async (resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    await client.sendRequests({
      payload: payload,
      requestType: RequestType.CreateFlightRequestType,
      byteArrayBufferNo: 1,
      totalByteArrayBuffers: 1,
      responseLost: true
    });
    resolve(1);
  });
}

// Simulate update Flight Price with request lost
export function updateFlightPriceRequestWithRequestLost(
  flightIdentifier: number,
  newPrice: string
) {
  const payload = marshal({
    FlightIdentifier: flightIdentifier,
    NewPrice: newPrice
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    console.log('Sending a packet, will resend in 5 sec...');
    setTimeout(async () => {
      await client.sendRequests({
        payload: payload,
        requestType: RequestType.UpdateFlightPriceRequestType,
        byteArrayBufferNo: 1,
        totalByteArrayBuffers: 1
      });
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
    NewPrice: newPrice
  });

  return new Promise((resolve, reject) => {
    const client = new UDPClient(ip, 8080);
    console.log('Sending a packet, will resend in 5 sec...');
    setTimeout(async () => {
      await client.sendRequests({
        payload: payload,
        requestType: RequestType.UpdateFlightPriceRequestType,
        byteArrayBufferNo: 1,
        totalByteArrayBuffers: 1,
        responseLost: true
      });
    }, 5000);
  });
}
