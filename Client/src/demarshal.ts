import { Buffer } from 'buffer';
import { StatusCode } from './interfaces';

export function demarshal (buffer: Buffer) {
    let dec = new TextDecoder();
    const bufferString = dec.decode(buffer);
    
    let bufferSplit = bufferString.split(/\r?\n/).filter(e => e);
    bufferSplit = bufferSplit.map(e =>  e.trim());

    // Get status code from buffer
    let statusCode = processStatusCode(bufferSplit[1]);

    if(statusCode === StatusCode.Success)
        processData(bufferSplit[2]);

    else if(statusCode === StatusCode.BusinessLogicGenericError)
        return "Generic Error"
    else if(statusCode === StatusCode.MarshallerError)
        return "Marshal Error"
    else if(statusCode === StatusCode.NoMatchForSourceAndDestination)
        return "No match for source and destination"
    else if(statusCode === StatusCode.NoSuchFlightIdentifier)
        return "No flight identifier"
    else if(statusCode === StatusCode.InsufficientNumberOfAvailableSeats)
        return "Insufficient number of available seats"
}

function processStatusCode(statusStr: string) {
    const statusCode = Number(statusStr.split(' ')[1]);

    return statusCode;
}

function processData(raw: string) {
    const data = raw.split(' ');
    console.log(data);
}

const responseQ1 = `{StatusCode:1,Data:{FlightIdentifiers:512345}}`
const responseQ2 = `{StatusCode:1,Data:{DepartureTime:1400,Airfare:1200.00,TotalAvailableSeats:30}}`
const responseQ4 = `{StatusCode:1,Data:{TotalAvailableSeats:30}}`
const responseUpdate = `{StatusCode:1,Data:{FlightIdentifier:2,SourceLocation:'Singapore',DestinationLocation:'Denmark',DepartureTime:1200,Airfare:1200.00,TotalAvailableSeats:30}}`
const responseCreate = `{StatusCode:1,Data:{FlightIdentifier:2}}`

const buffer = Buffer.from(responseQ1);
demarshal(buffer);