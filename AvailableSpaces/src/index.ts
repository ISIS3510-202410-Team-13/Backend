import express, { Express, Request, Response } from "express";
import dotenv from "dotenv";
import { MeetingRequest, AvailableSpace } from "./types";
import { findAvailableSpaces } from "./service";
import { validateBody } from "./validate";

dotenv.config();

const app: Express = express();
const port = Number(process.env.PORT) || 3000;
const host = process.env.HOST || "localhost";

app.use(express.json());
app.use(validateBody);

app.get("/", async (req: Request<{}, {}, MeetingRequest>, res: Response<AvailableSpace[]>) => {
  const { dayOfWeek, startTime, endTime } = req.body;
  const availableSpaces: AvailableSpace[] = await findAvailableSpaces(dayOfWeek, startTime, endTime);
  res.send(availableSpaces);
});

app.listen(port, host, () => {
  console.log(`[available-spaces]: Available Spaces Server is running at http://${host}:${port}`);
});
