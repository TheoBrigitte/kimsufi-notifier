"use client";

import { useState } from "react";
import useWebSocket from "react-use-websocket";
import ServersTable from "./components/serversTableComponent/ServersTableComponent";
import HeaderComponent from "./components/headerComponent/HeaderComponent";
import { Server, ErrorNull } from "./components/types";

export default function Home() {
  // Get the websocket URL from the environment variable
  const websocketURL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || "wss://localhost";
  // data is the list of servers received from the websocket server
  const [data, setData] = useState<Server[]>([]);
  // connectionRestored is used to display a message when the websocket connection is restored
  const [connectionRestored, setConnectionRestored] = useState(false);
  // lastMessage is the timestamp of the last message received from the websocket server
  const [lastMessage, setLastMessage] = useState(0);
  // error is the error message from the websocket connection
  const [error, setError] = useState<ErrorNull>(null);

  // Connect to the websocket server
  useWebSocket(websocketURL, {
    share: false,
    onMessage: (event) => {
      setError(null);
      setData(JSON.parse(event.data));
      setLastMessage(Date.now());
    },
    onOpen: () => {
      setError(null);
      setConnectionRestored(true);
    },
    onClose: () => setError(Error("socket closed")),
    onError: () => setError(Error("socket error")),
    // Reconnect with exponential backoff
    reconnectInterval: (attemptNumber) => Math.min(Math.pow(2, attemptNumber) * 1000, 10000),
    // Reconnect indefinitely
    shouldReconnect: () => true,
  });

  return (
    <div className="Home flex flex-col items-center pt-10 pb-20 px-10">
      <HeaderComponent
        error={error}
        data={data}
        connectionRestored={connectionRestored}
        setConnectionRestored={setConnectionRestored}
        lastMessage={lastMessage}
      />
      <ServersTable data={data} />
    </div>
  );
}
