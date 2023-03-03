export type PendingRequest = {
    data: string;
    attempts: number;
}

export enum StatusCode {
    Success = 1,
    BusinessLogicGenericError = 2,
    MarshallerError = 3,
    NoMatchForSourceAndDestination = 4,
    NoSuchFlightIdentifier = 5,
    InsufficientNumberOfAvailableSeats = 6
}

