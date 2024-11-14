'use client'

import React, { useState } from 'react';
import ServersTable from './components/server';
import useWebSocket  from 'react-use-websocket';
import { Server } from './components/types';

type ErrorNull = Error | null;

type StatusProps = {
  error: ErrorNull;
  data: Server[];
  connectionRestored: boolean;
  setConnectionRestored: (value: boolean) => void;
  lastMessage: number;
};

// Status components displays the status of the websocket connection
function Status({error, data, connectionRestored, setConnectionRestored, lastMessage}: StatusProps) {
  let message;
  let fadeOut;

  if (error) {
    // Display error message
    message =
      <>
        <div>Failed to load server list</div>
        <div className="text-orange-700">{error.toString()}</div>
      </>
  } else if (!data||data==undefined) {
    // No data yet, still loading
    message = <div>Loading ...</div>
  } else if (connectionRestored) {
    // Connection restored message
    message = <div className="text-green-700">Websocket connected</div>
    fadeOut = "transition-opacity duration-[2000ms] opacity-0"
    // Reset connectionRestored after 5 seconds
    setTimeout(function() {
      setConnectionRestored(false);
    }, 5000);
  } else if (lastMessage > 0) {
    // Last message received timestamp
    const date = new Date(lastMessage).toTimeString().split(' ')[0]
    message = <div>Last update received at {date}</div>
  }

  // Hide the message if there is no message to display
  const hidden = !message ? "hidden" : ""

  return (
    <div className={`basis-1/4 flex flex-col justify-center font-mono ${hidden} ${fadeOut}`}>
      {message}
    </div>
  )
}

export default function Home() {
  // Get the websocket URL from the environment variable
  const websocketURL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost';

  // Initialize the state variables
  //
  // data is the list of servers received from the websocket server
  const emptyData = Array<Server>();
  const [data, setData] = useState(emptyData);
  // connectionRestored is used to display a message when the websocket connection is restored
  const [connectionRestored, setConnectionRestored] = useState(false);
  // lastMessage is the timestamp of the last message received from the websocket server
  const [lastMessage, setLastMessage] = useState(0);
  // error is the error message from the websocket connection
  const [error, setError] = useState(null as ErrorNull);

  // Connect to the websocket server
  useWebSocket(
    websocketURL,
    {
      share: false,
      onMessage: (event) => {
        setError(null)
        setData(JSON.parse(event.data))
        setLastMessage(Date.now())
      },
      onOpen: () => {
        setError(null)
        setConnectionRestored(true)
      },
      onClose: () => setError(Error("socket closed")),
      onError: () => setError(Error("socket error")),
      // Reconnect with exponential backoff
      reconnectInterval: (attemptNumber) => Math.min(Math.pow(2, attemptNumber) * 1000, 10000),
      // Reconnect indefinitely
      shouldReconnect: () => true,
    },
  );

  return (
    <div className="flex flex-row justify-center">
      <div className="pt-10 pb-20 px-10">
        <div className="flex flex-row min-w-fit justify-center flex-nowrap text-nowrap">
          <div className="basis-1/4 flex-none"></div>
          <div className="basis-2/4 px-40 py-5 flex-none text-center text-xl font-bold">OVH Eco server availability</div>
          <Status
            error={error}
            data={data}
            connectionRestored={connectionRestored}
            setConnectionRestored={setConnectionRestored}
            lastMessage={lastMessage}
          />
        </div>
        <ServersTable data={data} />
      </div>
    </div>
  );
}
