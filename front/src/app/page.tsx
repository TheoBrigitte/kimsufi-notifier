'use client'

import React from 'react'
import ServersTable from './components/server';
import useSWRSubscription from 'swr/subscription';

function Status({ error, data, connectionRestored, setConnectionRestored, lastMessage, setLastMessage }) {
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
    fadeOut = "transition-opacity duration-[1000ms] opacity-0"
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
  const [data, setData] = React.useState(null);
  const [connectionRestored, setConnectionRestored] = React.useState(false);
  const [lastMessage, setLastMessage] = React.useState(0);

  const startWS = (key, { next }) => {
    let socket = new WebSocket("ws://127.0.0.1:9779/list", 'echo-protocol');
    socket.addEventListener('message', (event) => {
      const res = JSON.parse(event.data)
      next(null, res)
      setData(res)
      setLastMessage(Date.now())
    })
    socket.addEventListener('error', (event) => {
      next(event.error)
      //ws.close()
    })
    socket.addEventListener('close', (event) => {
      next(new Error("socket closed"))
      setTimeout(function() {
        startWS(key, { next })
      }, 5000);
    })
    socket.addEventListener('open', () => {
      setConnectionRestored(true)
    })
    return () => socket.close()
  }
  
  const { error } = useSWRSubscription('/list', startWS)
  const status = <Status
    error={error}
    data={data}
    connectionRestored={connectionRestored}
    setConnectionRestored={setConnectionRestored}
    lastMessage={lastMessage}
    setLastMessage={setLastMessage}
  />

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
        <div className="flex flex-row min-w-fit flex-nowrap text-nowrap">
          <div className="basis-1/4 flex-none"></div>
          <div className="basis-2/4 px-40 py-5 flex-none text-center text-xl font-bold">OVH Eco server availability</div>
          {status}
        </div>
        <ServersTable data={data} />
      </div>
    </div>
  );
}
