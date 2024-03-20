type Day = 'l' | 'm' | 'i' | 'j' | 'v' | 's' | 'd';

interface MeetingRequest {
  dayOfWeek: Day;
  startTime: string;
  endTime: string;
}

interface AvailableSpace {
  building: string;
  room: string;
  availableFrom: string;
  availableUntil: string;
  minutesAvailable: number;
}

interface UniandesCourseSection {
  schedules: UniandesCourseSchedule[]
}

type UniandesCourseSchedule = {
  [key in Day]: string | null;
} & {
  time_ini: string;
  time_fin: string;
  classroom: string;
}

type RoomReservations = {
  [key in Day]?: TimeBlock[];
} & {
  building: string;
  room: string;
};

interface TimeBlock {
  startMinute: number;
  endMinute: number;
}

export { Day, MeetingRequest, AvailableSpace, UniandesCourseSection, RoomReservations, TimeBlock};
