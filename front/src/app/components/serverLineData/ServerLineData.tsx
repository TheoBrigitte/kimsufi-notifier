"use client";

import React, { useState, useEffect } from "react";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle, faCheckCircle, faTimesCircle } from "@fortawesome/free-solid-svg-icons";

import { Server, Status } from "../types";

interface Props {
  server: Server;
}

// Define the color and icon the status column
const statusColor = new Map<string, Status>([
  ["available", { color: "text-lime-600", icon: faCheckCircle }],
  ["unavailable", { color: "text-rose-600", icon: faTimesCircle }],
]);

export function ServerLineData({ server }: Props) {
  const [prevProps, setPrevProps] = useState<Server>(server);
  const [differences, setDifferences] = useState<boolean>(false);

  const dataEquality = JSON.stringify(prevProps) === JSON.stringify(server);

  useEffect(() => {
    if (!dataEquality) {
      console.log("change detected");
      setDifferences(true);
      setPrevProps(server);
    }
  }, [dataEquality, server]);

  useEffect(() => {
    const timerHighlight = setTimeout(() => setDifferences(false), 1000);
    return () => clearTimeout(timerHighlight);
  }, [differences]);

  return (
    <tr
      className={`font-mono text-sm font-medium ${differences ? "bg-yellow-200" : "transition duration-1000 delay-150 shadow-md rounded-xl bg-slate-100"}`}
    >
      <td className="rounded-l-xl">{server.name}</td>
      <td>{server.cpu}</td>
      <td>{server.memory}</td>
      <td>{server.storage}</td>
      <td>{server.bandwidth}</td>
      <td>{`${server.price} ${server.currencyCode}`}</td>
      <td className="flex justify-end py-6">
        <div>{server.status}</div>
        <div className={statusColor.get(server.status)?.color + " pl-2"}>
          <FontAwesomeIcon icon={statusColor.get(server.status)?.icon || faCircle} />
        </div>
      </td>
      <td className="rounded-r-xl">{server.datacenters?.join(", ") || "-"}</td>
    </tr>
  );
}
