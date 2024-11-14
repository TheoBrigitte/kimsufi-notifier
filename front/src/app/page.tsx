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

function Status({error, data, connectionRestored, setConnectionRestored, lastMessage}: StatusProps) {
  let message;
  let fadeOut;

  if (error) {
    message =
      <>
        <div>Failed to load server list</div>
        <div className="text-orange-700">{error.toString()}</div>
      </>
  } else if (!data||data==undefined) {
    message = <div>Loading ...</div>
  } else if (connectionRestored) {
    message = <div className="text-green-700">Websocket connected</div>
    fadeOut = "transition-opacity duration-[2000ms] opacity-0"
    setTimeout(function() {
      setConnectionRestored(false);
    }, 5000);
  } else if (lastMessage > 0) {
    const date = new Date(lastMessage).toTimeString().split(' ')[0]
    message = <div>Last update received at {date}</div>
  }

  const hidden = !message ? "hidden" : ""

  return (
    <div className={`basis-1/4 flex flex-col justify-center font-mono ${hidden} ${fadeOut}`}>
      {message}
    </div>
  )
}

export default function Home() {
  const websocketURL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost';

  const emptyData = Array<Server>();
  const [data, setData] = useState(emptyData);
  const [connectionRestored, setConnectionRestored] = useState(false);
  const [lastMessage, setLastMessage] = useState(0);
  const [error, setError] = useState(null as ErrorNull);

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
      reconnectInterval: (attemptNumber) => Math.min(Math.pow(2, attemptNumber) * 1000, 10000),
      shouldReconnect: () => true,
    },
  );

  //const setServers = (data) => {
  //  setData(data);
  //  return 5000;
  //}
  //const opts = {
  //  revalidateOnFocus: false,
  //  refreshInterval: setServers
  //}
  //const { error, isLoading } = useSWR('https://127.0.0.1:8080/list', getServers, opts);

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
