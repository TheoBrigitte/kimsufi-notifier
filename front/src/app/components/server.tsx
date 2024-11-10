'use client';

import { ReactNode } from "react";
import React from 'react';
import { useEffect, useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle, faTimesCircle } from '@fortawesome/free-regular-svg-icons';

interface Props {
    data: Object
}

function ServerLines({ category, servers }) {
  let statusColor = new Map<string, Object>([
      ["available", {"color":"text-lime-600", "icon":faCheckCircle}],
      ["unavailable", {"color":"text-rose-600", "icon":faTimesCircle}],
  ]);

  const [rowsData, setRowsData] = useState({});
  const [changedRows, setChangedRows] = useState(new Set());
  useEffect(() => {
    const newRowsData = {};

    // Detect changes by comparing the current data to the previous data
    servers.forEach((server) => {
      const previousData = rowsData[server.planCode];
      const currentData = server;

      // Check if the data has changed
      const equal = JSON.stringify(previousData) === JSON.stringify(currentData)
      if (previousData && !equal) {
        setChangedRows((prevChangedRows) => new Set(prevChangedRows).add(server.planCode));
      }

      // Update the data for the current row
      newRowsData[server.planCode] = currentData;
    });

    setRowsData(newRowsData);
  }, [servers]); // Re-run effect when servers data changes

    // Clear highlight after a delay
  useEffect(() => {
    if (changedRows.size > 0) {
      const timer = setTimeout(() => {
        setChangedRows(new Set()); // Clear changed rows
      }, 0); // Duration for highlight in milliseconds

      return () => clearTimeout(timer);
    }
  }, [changedRows]);

  return (
      servers.map((server) => (
      <tr ref={React.createRef()} key={category + server.planCode} className={`${changedRows.has(server.planCode) ? 'bg-yellow-200' : 'transition duration-1000 delay-150 even:bg-blue-300 odd:bg-blue-100'} font-mono`}>
        <td>{server.name}</td>
        <td>{server.cpu}</td>
        <td>{server.memory}</td>
        <td>{server.storage}</td>
        <td>{server.bandwidth}</td>
        <td>{server.price} {server.currencyCode}</td>
        <td className="flex flex-row justify-end">{server.status}<div className={statusColor.get(server.status).color + " pl-2"}><FontAwesomeIcon icon={statusColor.get(server.status).icon} /></div></td>
        <td>{server.datacenters?.join(", ")||"-"}</td>
      </tr>
    ))
  )
}

const ServersTable = ({data} : Props) => {
  if (!data) {
    return
  }

  const serversByCategory = Object.groupBy(data, ({ category }) => category);

  let categoryOrder = { "Kimsufi": [], "So you Start": [], "Rise": [], "": [], }
  const ordered = Object.assign(categoryOrder, serversByCategory);

  const tableBody = Object.entries(ordered).map(([category, servers]) => (
    <>
      <tr key={category + " name"}><td className="p-2 font-mono" colSpan={6}>{category||"Uncategorized"}</td></tr>
      <ServerLines key={category + " servers"} catagory={category} servers={servers} />
    </>
  ));

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
          {tableBody}
        </tbody>
      </table>
  );
};

export default ServersTable;
