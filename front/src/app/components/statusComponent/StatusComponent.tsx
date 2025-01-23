import { ReactNode } from "react";
import { StatusProps } from "./type";

// Status components displays the status of the websocket connection
export default function Status({
  error,
  data,
  connectionRestored,
  setConnectionRestored,
  lastMessage,
}: StatusProps) {
  let statusDisplay: ReactNode | null = null;
  const fadeOut: string = "transition-opacity duration-[2000ms] opacity-0";

  const errorDisplay = (error: Error) => (
    <>
      <div>Failed to load server list</div>
      <div className="text-orange-700">{error.toString()}</div>
    </>
  );

  const loaderDisplay: ReactNode = <div>Loading ...</div>;

  const restoredDisplay: ReactNode = (
    <div className="text-green-700">Websocket connected</div>
  );

  const lastMessageDisplay: ReactNode = (
    <div>
      Last update received at{" "}
      {new Date(lastMessage).toTimeString().split(" ")[0]}
    </div>
  );

  // Display error message
  if (error) statusDisplay = errorDisplay(error);
  // Display loader
  if (!data.some((e) => e)) statusDisplay = loaderDisplay;
  // Connection restored message
  if (connectionRestored) statusDisplay = restoredDisplay;
  // Reset connectionRestored after 5 seconds
  setTimeout(() => setConnectionRestored(false), 5000);
  // Last message received timestamp
  if (lastMessage > 0) statusDisplay = lastMessageDisplay;

  return (
    <div
      className={`
        basis-1/4 flex flex-col justify-center font-mono
        ${!statusDisplay ? "hidden" : ""}
        ${connectionRestored ? fadeOut : ""}
      `}
    >
      {statusDisplay}
    </div>
  );
}
