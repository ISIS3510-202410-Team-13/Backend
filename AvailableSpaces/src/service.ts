import { Day } from "./types";

const findAvailableSpaces = async (dayOfWeek: Day, startTime: string, endTime: string) => {
    return [
        {
            building: 'A',
            room: '1',
            availableFrom: "0930",
            availableUntil: "1930"
        },
        {
            building: 'A',
            room: '2',
            availableFrom: "0930",
            availableUntil: "1930"
        }
    ]
};

export { findAvailableSpaces };
