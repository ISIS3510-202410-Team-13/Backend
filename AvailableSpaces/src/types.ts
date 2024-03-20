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
}

export { Day, MeetingRequest, AvailableSpace };
