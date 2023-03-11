import * as readline from "readline/promises";
import { stdin as input, stdout as output } from "process";
import {
  createSeatReservationRequest,
  getFlightIdentifier,
  getFlightInformation,
  monitorSeatUpdatesCallbackRequest,
  updateFlightPriceRequest,
  createFlightRequest,
} from "./stubs";

const rl = readline.createInterface({ input, output });

export async function userInterface() {
  let option = 0;
  console.log("-----Hello, Welcome to Flight Information System-----");
  console.log("1. Query Flight Identifier");
  console.log("2. Query Flight Information");
  console.log("3. Make a Flight Seat Reservation");
  console.log("4. Monitor Flight Seats Information");
  console.log("5. Update The Flight Ticket Price");
  console.log("6. Create a Flight Request");
  option = Number(await rl.question("What do you wish to do?\n"));
  while (option < 1 || option > 6) {
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
  }

  return inputs;
}

async function q1() {
  const sourceLocation = await rl.question("Input your Source Location\n");
  const destinationLocation = await rl.question(
    "Input your destination location\n"
  );

  getFlightIdentifier(sourceLocation, destinationLocation);
}

async function q2() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );

  getFlightInformation(flightIdentifier);
}

async function q3() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );
  const seatsToReserve = Number(
    await rl.question("Input your Seats to Reserve\n")
  );

  createSeatReservationRequest(flightIdentifier, seatsToReserve);
}

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

async function q5() {
  const flightIdentifier = Number(
    await rl.question("Input your Flight Identifier Number\n")
  );
  const newPrice = Number(await rl.question("Input your new price\n"));

  updateFlightPriceRequest(flightIdentifier, newPrice);
}

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

userInterface();
