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
  SourceLocation: string;
  DestinationLocation: string;
};

export type GetFlightIdentifiersResponse = {
  FlightIdentifiers: number[];
};

export type GetFlightInformationRequest = {
  FlightIdentifier: number;
};

export type GetFlightInformationResponse = {
  DepartureTime: number;
  Airfare: number;
  TotalAvailableSeats: number;
};

export type MakeSeatReservationRequest = {
  FlightIdentifier: number;
  SeatsToReserve: number;
};

export type MonitorSeatUpdatesCallbackRequest = {
  FlightIdentifier: number;
  LengthOfMonitorIntervalInSeconds: number;
};

export type MonitorSeatUpdatesCallbackResponse = {
  TotalAvailableSeats: number;
};

export type UpdateFlightPriceRequest = {
  FlightIdentifier: number;
  NewPrice: number;
};

export type UpdateFlightPriceResponse = {
  FlightIdentifier: number;
  SourceLocation: string;
  DestinationLocation: string;
  DepartureTime: bigint;
  Airfare: number;
  TotalAvailableSeats: number;
};

export type CreateFlightRequest = {
  SourceLocation: string;
  DestinationLocation: string;
  DepartureTime: bigint;
  Airfare: number;
  TotalAvailableSeats: number;
};

export type CreateFlightResponse = {
  FlightIdentifier: number;
};
