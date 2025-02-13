import React, { useState, useEffect } from "react";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle, faCheckCircle, faTimesCircle } from "@fortawesome/free-solid-svg-icons";

import { Server, Status } from "../types";

export function ServerLineData({ server }: { server: Server }) {
  const [prevProps, setPrevProps] = useState<Server>(server);
  const [differences, setDifferences] = useState<boolean | null>(null);

  const dataEquality = JSON.stringify(prevProps) === JSON.stringify(server);

  useEffect(() => {
    console.log("start useEffect");
    if (!dataEquality) {
      setDifferences(true);
      setPrevProps(server);
      // Réinitialisation après 0.5s
      // const timer = setTimeout(() => setDifferences(null), 500);
      // return () => clearTimeout(timer);
    }
  }, [dataEquality, server]);

  // Define the color and icon the status column
  const statusColor = new Map<string, Status>([
    ["available", { color: "text-lime-600", icon: faCheckCircle }],
    ["unavailable", { color: "text-rose-600", icon: faTimesCircle }],
  ]);

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
