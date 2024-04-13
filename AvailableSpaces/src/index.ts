import express, { Express, Request, Response } from "express";
import { MeetingRequest, AvailableSpace } from "./types";
import { findAvailableSpaces } from "./service";
import { validateBody } from "./validate";

const app: Express = express();
const port = process.env.PORT || 3000;

app.use(express.json());

app.get("/health", (req: Request, res: Response) => {
  res.status(200).send("Available Spaces Server is running").end();
})

app.post("/spaces", validateBody)
app.post("/spaces", async (req: Request<{}, {}, MeetingRequest>, res: Response<AvailableSpace[]>) => {
  const { dayOfWeek, startTime, endTime } = req.body;
  const availableSpaces: AvailableSpace[] = await findAvailableSpaces(dayOfWeek, startTime, endTime);
  res.status(200).send(availableSpaces).end();
});

app.listen(port, () => {
  console.log(`[available-spaces] 🚀 Available Spaces Server is running at https://localhost:${port}`);
});
