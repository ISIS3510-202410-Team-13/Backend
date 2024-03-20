import { Day, UniandesCourseSection, RoomReservations, TimeBlock, AvailableSpace } from "./types";
const groupBy = require('lodash/groupBy');
import axios from 'axios';

const buildingsWhiteList = ["ML", "SD", "W", "R", "O", "C", "LL", "B", "RGD", "AU"];
let roomsReservations: RoomReservations[] = [];

const transformToMinutes = (time: string) => {
  const hours = parseInt(time.slice(0, 2));
  const minutes = parseInt(time.slice(2));
  return hours * 60 + minutes;
}

const transformToTime = (minutes: number) => {
  const hours = Math.floor(minutes / 60);
  const minutesLeft = minutes % 60;
  return `${hours.toString().padStart(2, "0")}${minutesLeft.toString().padStart(2, "0")}`;
}

const sectionToReservation = (course: UniandesCourseSection): (RoomReservations | undefined)[] => {
  return course.schedules.map(schedule => {
    try {
      const building = schedule.classroom.split("_")[0].replace(/[^a-zA-Z0-9]/g, '');
      const room = schedule.classroom.split("_")[1].replace(/[^a-zA-Z0-9]/g, '');
      const startMinute = transformToMinutes(schedule.time_ini);
      const endMinute = transformToMinutes(schedule.time_fin);
      const days = (["l", "m", "i", "j", "v", "s", "d"] as Day[]).filter(day => schedule[day]);
      const daysWithTime = days.map(day => ({ day : { startMinute, endMinute } }));
      return { building, room, ...daysWithTime };
    } catch {
      return undefined;
    }
  }).filter(reservation => !!reservation);
}

const mergeReservationsByRoom = (reservations: RoomReservations[]): RoomReservations[] => {
  const groupedReservations: {[key in string]: RoomReservations[]} = groupBy(reservations, (reservation: RoomReservations) => `${reservation.building}-${reservation.room}`);
  return Object.values(groupedReservations).map((reservations: RoomReservations[]) => {
    const building = reservations[0].building;
    const room = reservations[0].room;
    const days = reservations.reduce((acc, reservation) => {
      (["l", "m", "i", "j", "v", "s", "d"] as Day[]).forEach((key) => {
        acc[key] = [...acc[key], ...(reservation[key] || [])];
      });
      return acc;
    }, Object.fromEntries(["l", "m", "i", "j", "v", "s", "d"].map(day => [day, [] as TimeBlock[] ])));
    const daysWithMergedTimeBlocks = Object.fromEntries(Object.entries(days).map(([day, timeBlocks]) => [day, joinOverlappingTimeBlocks(timeBlocks)]));
    return { building, room, ...daysWithMergedTimeBlocks };
  });
}

const joinOverlappingTimeBlocks = (timeBlocks: TimeBlock[]): TimeBlock[] => {
  const sortedTimeBlocks = timeBlocks.sort((a, b) => a.startMinute - b.startMinute);
  return sortedTimeBlocks.reduce((acc, timeBlock) => {
    if (acc.length === 0) return [timeBlock];
    const lastTimeBlock = acc[acc.length - 1];
    if (timeBlock.startMinute <= lastTimeBlock.endMinute) {
      lastTimeBlock.endMinute = Math.max(lastTimeBlock.endMinute, timeBlock.endMinute);
    } else {
      acc.push(timeBlock);
    }
    return acc;
  }, [] as TimeBlock[]);
};

const fetchUniandesAPI = () => {
  axios.get('https://ofertadecursos.uniandes.edu.co/api/courses')
    .then(response => response.data)
    .then((data: UniandesCourseSection[]) => data.flatMap(sectionToReservation))
    .then(reservations => reservations.filter(reservation => buildingsWhiteList.includes(reservation!.building)))
    .then(reservations => roomsReservations = mergeReservationsByRoom(reservations as RoomReservations[]))
    .then(() => console.log('[available-spaces] ✅ Updated reservations at ' + new Date().toLocaleString('es-CO', { timeZone: 'America/Bogota' })))
    .catch(error => console.error('[available-spaces] 🚨', error));
}


fetchUniandesAPI();
setInterval(fetchUniandesAPI, 1000 * 60 * 24);

const checkIfRoomIsAvailable = (roomReservations: RoomReservations, dayOfWeek: Day, startMinute: number, endMinute: number) => {
  if(roomReservations.building === "O" && roomReservations.room === "104") console.log(roomReservations, dayOfWeek, startMinute, endMinute);
  const dayReservations = roomReservations[dayOfWeek];
  if (!dayReservations || dayReservations.length === 0) return true;
  return !dayReservations.some(reservation => {
    const reservationStartIsInside = startMinute <= reservation.startMinute && endMinute > reservation.startMinute;
    const reservationEndIsInside = startMinute < reservation.endMinute && endMinute >= reservation.endMinute;
    return reservationStartIsInside || reservationEndIsInside;
  });
};

const calculateAvailableSpace = (roomReservations: RoomReservations, dayOfWeek: Day, startMinute: number, endMinute: number): AvailableSpace => {
  const dayReservations = roomReservations[dayOfWeek];
  if (!dayReservations || dayReservations.length === 0) return { building: roomReservations.building, room: roomReservations.room, availableFrom: "0000", availableUntil: "2359", minutesAvailable: 24 * 60 };
  const startAvailable = dayReservations.reduce((acc, reservation) => {
    if (startMinute >= reservation.endMinute) {
      return reservation.endMinute;
    }
    return acc;
  }, 0);
  const endAvailable = dayReservations.reverse().reduce((acc, reservation) => {
    if (endMinute <= reservation.startMinute) {
      return reservation.startMinute;
    }
    return acc;
  }, 24 * 60);
  return { building: roomReservations.building, room: roomReservations.room, availableFrom: transformToTime(startAvailable), availableUntil: transformToTime(endAvailable), minutesAvailable: endAvailable - startAvailable};
}

const findAvailableSpaces = async (dayOfWeek: Day, startTime: string, endTime: string): Promise<AvailableSpace[]> => {
  const startMinute = transformToMinutes(startTime);
  const endMinute = transformToMinutes(endTime);
  const availableRooms = roomsReservations.filter(reservation => checkIfRoomIsAvailable(reservation, dayOfWeek, startMinute, endMinute));
  return availableRooms.map(room => calculateAvailableSpace(room, dayOfWeek, startMinute, endMinute));
};

export { findAvailableSpaces };
