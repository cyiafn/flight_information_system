export type PendingRequest = {
  data: string;
  attempts: number;
};

export enum StatusCode {
  Success = 1,
  BusinessLogicGenericError = 2,
  MarshallerError = 3,
  NoMatchForSourceAndDestination = 4,
  NoSuchFlightIdentifier = 5,
  InsufficientNumberOfAvailableSeats = 6,
}

export enum RequestType {
  PingRequestType = 1,
  GetFlightIdentifiersRequestType = 2,
  GetFlightInformationRequestType = 3,
  MakeSeatReservationRequestType = 4,
  MonitorSeatUpdatesRequestType = 5,
  UpdateFlightPriceRequestType = 6,
  CreateFlightRequestType = 7,
}

export enum ResponseType {
  PingResponseType = 101,
  GetFlightIdentifiersResponseType = 102,
  GetFlightInformationResponseType = 103,
  MakeSeatReservationResponseType = 104,
  MonitorSeatUpdatesResponseType = 105,
  UpdateFlightPriceResponseType = 106,
  CreateFlightResponseType = 107,
  MonitorSeatUpdatesCallbackType = 201,
}

export type GetFlightIdentifiersRequest = {
  Discriminator: RequestType.GetFlightIdentifiersRequestType;
  SourceLocation: string;
  DestinationLocation: string;
};

export function instanceOfGetFlightIdentifiersRequest(
  obj: any
): obj is GetFlightIdentifiersRequest {
  return obj.Discriminator === RequestType.GetFlightIdentifiersRequestType;
}

export type GetFlightIdentifiersResponse = {
  FlightIdentifiers: number[];
};

export type GetFlightInformationRequest = {
  Discriminator: RequestType.GetFlightInformationRequestType;
  FlightIdentifier: number;
};

export function instanceOfGetFlightInformationRequest(
  obj: any
): obj is GetFlightInformationRequest {
  return obj.Discriminator === RequestType.GetFlightInformationRequestType;
}

export type GetFlightInformationResponse = {
  DepartureTime: number;
  Airfare: number;
  TotalAvailableSeats: number;
};

export type MakeSeatReservationRequest = {
  Discriminator: RequestType.MakeSeatReservationRequestType;
  FlightIdentifier: number;
  SeatsToReserve: number;
};

export function instanceOfMakeSeatReservationRequest(
  obj: any
): obj is MakeSeatReservationRequest {
  return obj.Discriminator === RequestType.MakeSeatReservationRequestType;
}

export type MonitorSeatUpdatesCallbackRequest = {
  Discriminator: RequestType.MonitorSeatUpdatesRequestType;
  FlightIdentifier: number;
  LengthOfMonitorIntervalInSeconds: number;
};

export function instanceOfMonitorSeatUpdatesRequest(
  obj: any
): obj is MonitorSeatUpdatesCallbackRequest {
  return obj.Discriminator === RequestType.MonitorSeatUpdatesRequestType;
}

export type MonitorSeatUpdatesCallbackResponse = {
  TotalAvailableSeats: number;
};

export type UpdateFlightPriceRequest = {
  Discriminator: RequestType.UpdateFlightPriceRequestType;
  FlightIdentifier: number;
  NewPrice: number;
};

export function instanceOfUpdateFlightPriceRequest(
  obj: any
): obj is UpdateFlightPriceRequest {
  return obj.Discriminator === RequestType.UpdateFlightPriceRequestType;
}

export type UpdateFlightPriceResponse = {
  FlightIdentifier: number;
  SourceLocation: string;
  DestinationLocation: string;
  DepartureTime: number;
  Airfare: number;
  TotalAvailableSeats: number;
};

export type CreateFlightRequest = {
  Discriminator: RequestType.CreateFlightRequestType;
  SourceLocation: string;
  DestinationLocation: string;
  DepartureTime: number;
  Airfare: number;
  TotalAvailableSeats: number;
};

export function instanceOfCreateFlightRequest(
  obj: any
): obj is CreateFlightRequest {
  return obj.Discriminator === RequestType.CreateFlightRequestType;
}

export type CreateFlightResponse = {
  FlightIdentifier: number;
};
