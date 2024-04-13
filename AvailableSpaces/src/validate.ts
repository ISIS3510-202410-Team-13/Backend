import { Request, Response, NextFunction } from "express";

const validateBody = (req: Request, res: Response, next: NextFunction) => {
  
  if (!req.body) {
    res.status(400).send({ message: 'Request body is missing' }).end();
    return;
  }

  const { dayOfWeek, startTime, endTime } = req.body;

  if (dayOfWeek === undefined || startTime === undefined || endTime === undefined) {
    res.status(400).send({ message: 'Missing required fields in request body' }).end();
    return;
  }

  if (typeof dayOfWeek !== "string" || typeof startTime !== "string" || typeof endTime !== "string") {
    res.status(400).send({ message: 'Invalid data types in request body' }).end();
    return;
  }

  if (dayOfWeek.length !== 1 || !["l", "m", "i", "j", "v", "s", "d"].includes(dayOfWeek)) {
    res.status(400).send({ message: 'Invalid day of week' }).end();
    return;
  }

  if (!/^\d{4}$/.test(startTime) || !/^\d{4}$/.test(endTime)) {
    res.status(400).send({ message: 'Invalid time format' }).end();
    return;
  }

  if (startTime >= endTime) {
    res.status(400).send({ message: 'Start time must be before end time' }).end();
    return;
  }

  if (startTime < "0000" || startTime >= "2400" || endTime < "0000" || endTime >= "2400") {
    res.status(400).send({ message: 'Invalid time range' }).end();
    return;
  }

  if (startTime.slice(-2) > "59" || endTime.slice(-2) > "59") {
    res.status(400).send({ message: 'Invalid minutes' }).end();
    return;
  }

  next();
}

export { validateBody };
