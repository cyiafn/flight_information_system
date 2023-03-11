import * as readline from "readline/promises";
import { stdin as input, stdout as output } from "process";
import {
  createSeatReservationRequest,
  getFlightIdentifier,
  getFlightInformation,
  monitorSeatUpdatesCallbackRequest,
  updateFlightPriceRequest,
  createFlightRequest,
  createFlightWithRequestLost,
  createFlightWithResponseLost,
} from "./stubs";

const rl = readline.createInterface({ input, output });

// Text-based UI to select options
export async function userInterface() {
  let option = 0;
  console.log("-----Hello, Welcome to Flight Information System-----");
  console.log("1. Query Flight Identifier");
  console.log("2. Query Flight Information");
  console.log("3. Make a Flight Seat Reservation");
  console.log("4. Monitor Flight Seats Information");
  console.log("5. Update The Flight Ticket Price");
  console.log("6. Create a Flight Request");
  console.log("7. Create a Flight Request with Request Lost");
  console.log("8. Create a Flight Request with Response Lost");
  console.log("Type q to quit.");
  option = Number(await rl.question("What do you wish to do?\n"));
  while (option < 1 || option > 8) {
    option = Number(await rl.question("Wrong Input\n"));
  }

  let inputs;
  switch (option) {
    case 1:
      inputs = await q1();
      break;
    case 2:
      inputs = await q2();
      break;
    case 3:
      inputs = await q3();
      break;
    case 4:
      inputs = await q4();
      break;
    case 5:
      inputs = await q5();
      break;
    case 6:
      inputs = await q6();
      break;
    case 7:
      inputs = await q7();
      break;
    case 8:
      inputs = await q8();
      break;
  }
  console.log(option);
  return String(option);
}

// Get flight identifier based on source location and destination location
async function q1() {
  const sourceLocation = await rl.question("Input your Source Location\n");
  const destinationLocation = await rl.question(
    "Input your destination location\n"
  );

  await getFlightIdentifier(sourceLocation, destinationLocation);
}

// Get Flight Information from the respective flight identifier number.
async function q2() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );

  getFlightInformation(flightIdentifier);
}

// Create seat reservation(s) based on flight identifier number and
// the number of seats to be reserved.
async function q3() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );
  const seatsToReserve = Number(
    await rl.question("Input your Seats to Reserve\n")
  );

  createSeatReservationRequest(flightIdentifier, seatsToReserve);
}

// Listen for seat number changes
async function q4() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );
  const lengthOfMonitorIntervalInSeconds = Number(
    await rl.question("Input your monitor interval in seconds\n")
  );

  monitorSeatUpdatesCallbackRequest(
    flightIdentifier,
    BigInt(lengthOfMonitorIntervalInSeconds)
  );
}

// Update flight price based on exisitng flight identifier
async function q5() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );
  const newPrice = Number(await rl.question("Input your new price\n"));

  updateFlightPriceRequest(flightIdentifier, newPrice);
}

// Create a flight request based on the source location, destination location,
// departure time, airfare and total available seats
async function q6() {
  const sourceLocation = await rl.question("Input your Source Location\n");
  const destinationLocation = await rl.question(
    "Input your destination location\n"
  );
  const departureTime = Number(
    await rl.question("Input your Departure Time\n")
  );
  const airfare = Number(
    await rl.question("Input the Airfare of your Flight\n")
  );
  const totalAvailableSeats = Number(
    await rl.question("Input the Total Available Seats of your Flight\n")
  );

  const dto = {
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
    DepartureTime: BigInt(departureTime),
    Airfare: airfare,
    TotalAvailableSeats: totalAvailableSeats,
  };

  createFlightRequest(dto);
}

// Create a flight request based on the source location, destination location,
// WITH REQUEST LOST
async function q7() {
  const sourceLocation = await rl.question("Input your Source Location\n");
  const destinationLocation = await rl.question(
    "Input your destination location\n"
  );
  const departureTime = Number(
    await rl.question("Input your Departure Time\n")
  );
  const airfare = Number(
    await rl.question("Input the Airfare of your Flight\n")
  );
  const totalAvailableSeats = Number(
    await rl.question("Input the Total Available Seats of your Flight\n")
  );

  const dto = {
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
    DepartureTime: BigInt(departureTime),
    Airfare: airfare,
    TotalAvailableSeats: totalAvailableSeats,
  };

  createFlightWithRequestLost(dto);
}

// Create a flight request based on the source location, destination location,
// WITH RESPONSE LOST
async function q8() {
  const sourceLocation = await rl.question("Input your Source Location\n");
  const destinationLocation = await rl.question(
    "Input your destination location\n"
  );
  const departureTime = Number(
    await rl.question("Input your Departure Time\n")
  );
  const airfare = Number(
    await rl.question("Input the Airfare of your Flight\n")
  );
  const totalAvailableSeats = Number(
    await rl.question("Input the Total Available Seats of your Flight\n")
  );

  const dto = {
    SourceLocation: sourceLocation,
    DestinationLocation: destinationLocation,
    DepartureTime: BigInt(departureTime),
    Airfare: airfare,
    TotalAvailableSeats: totalAvailableSeats,
  };

  createFlightWithResponseLost(dto);
}

async function main() {
  let input = "";
  while (input !== "q") {
    input = await userInterface();
  }
}

// Start of the program
main();
