"use client";

import React from "react";
import { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheckCircle, faTimesCircle, faCircle } from "@fortawesome/free-regular-svg-icons";
import { Server, Status } from "./types";

interface Props {
  data: Server[];
}

// ServerLines displays the servers table lines
function ServerLines({ category, servers }: { category: string; servers: Server[] }) {
  // Define the color and icon the status column
  const statusColor = new Map<string, Status>([
    ["available", { color: "text-lime-600", icon: faCheckCircle }],
    ["unavailable", { color: "text-rose-600", icon: faTimesCircle }],
  ]);

  // Highlight changed rows
  const [rowsData, setRowsData] = useState(new Map<string, Server>());
  const [changedRows, setChangedRows] = useState(new Set());
  useEffect(() => {
    const newRowsData = new Map<string, Server>();

    // Detect changes by comparing the current data to the previous data
    servers.forEach((server) => {
      const previousData = rowsData.get(server.planCode);
      const currentData = server;

      // Check if the data has changed
      const equal = JSON.stringify(previousData) === JSON.stringify(currentData);
      if (previousData && !equal) {
        setChangedRows((prevChangedRows) => new Set(prevChangedRows).add(server.planCode));
      }

      // Update the data for the current row
      newRowsData.set(server.planCode, currentData);
    });

    setRowsData(newRowsData);
  }, [servers]); // Re-run effect when servers data changes

  // Clear highlight after a delay
  useEffect(() => {
    if (changedRows.size > 0) {
      const timer = setTimeout(() => {
        setChangedRows(new Set()); // Clear changed rows
      }, 500); // Duration for highlight in milliseconds

      return () => clearTimeout(timer);
    }
  }, [changedRows]);

  return servers.map((server) => (
    <tr
      ref={React.createRef()}
      key={category + server.planCode}
      className={`font-mono text-sm font-medium ${changedRows.has(server.planCode) ? "bg-yellow-200" : "transition duration-1000 delay-150 even:bg-blue-300 odd:bg-blue-100"}`}
    >
      <td className="py-2">{server.name}</td>
      <td className="py-2">{server.cpu}</td>
      <td className="py-2">{server.memory}</td>
      <td className="py-2">{server.storage}</td>
      <td className="py-2">{server.bandwidth}</td>
      <td className="py-2">
        {server.price} {server.currencyCode}
      </td>
      <td className="flex flex-row justify-end py-2">
        {server.status}
        <div className={statusColor.get(server.status)?.color + " pl-2"}>
          <FontAwesomeIcon icon={statusColor.get(server.status)?.icon || faCircle} />
        </div>
      </td>
      <td>{server.datacenters?.join(", ") || "-"}</td>
    </tr>
  ));
}

// ServerCategories displays the servers table lines with a category line before
function ServerCategories({ ordered }: { ordered: { [key: string]: Server[] } }) {
  return Object.entries(ordered).map(([category, servers]) => (
    <React.Fragment key={category}>
      <tr>
        <td className="p-4 font-mono bg-gray-300 border-b-2 border-gray-500" colSpan={8}>
          {category}
        </td>
      </tr>
      <ServerLines category={category} servers={servers} />
    </React.Fragment>
  ));
}

const ServersTable = ({ data }: Props) => {
  if (!data || data.length === 0) return;
  // Group servers by category
  const serversByCategory = Object.groupBy(data, (server: Server) => server.category);
  // Define the order of categories
  const categoryOrder: { [key: string]: Server[] } = {
    Kimsufi: [],
    "So you Start": [],
    Rise: [],
    uncategorized: [],
  };
  // Merge the server and respect the categories order
  const ordered = Object.assign(categoryOrder, serversByCategory);

  return (
    <table className="text-nowrap">
      <thead>
        <tr>
          <th className="p-4">Name</th>
          <th className="p-4">CPU</th>
          <th className="p-4">RAM</th>
          <th className="p-4">Storage</th>
          <th className="p-4">Bandwidth</th>
          <th className="p-4">Price</th>
          <th className="p-4">Status</th>
          <th className="p-4">Datacenters</th>
        </tr>
      </thead>
      <tbody>
        <ServerCategories ordered={ordered} />
      </tbody>
    </table>
  );
};

export default ServersTable;
