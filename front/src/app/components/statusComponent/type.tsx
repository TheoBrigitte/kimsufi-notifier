import { ErrorNull, Server } from "../types";

type StatusProps = {
  error: ErrorNull;
  data: Server[];
  connectionRestored: boolean;
  setConnectionRestored: (value: boolean) => void;
  lastMessage: number;
};

export type { StatusProps };
