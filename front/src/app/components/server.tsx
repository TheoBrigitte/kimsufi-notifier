'use client';

import React from 'react';
import { useEffect, useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle, faTimesCircle, faCircle } from '@fortawesome/free-regular-svg-icons';
import { Server, Status } from './types';

interface Props {
    data: Server[];
}

function ServerLines({category, servers} : {category: string, servers: Server[]}) {
  const statusColor = new Map<string, Status>([
      ["available", {color:"text-lime-600", icon:faCheckCircle}],
      ["unavailable", {color:"text-rose-600", icon:faTimesCircle}],
  ]);

  const [rowsData, setRowsData] = useState(new Map<string, Server>());
  const [changedRows, setChangedRows] = useState(new Set());
  useEffect(() => {
    const newRowsData = new Map<string, Server>();

    // Detect changes by comparing the current data to the previous data
    servers.forEach((server) => {
      const previousData = rowsData.get(server.planCode);
      const currentData = server;

      // Check if the data has changed
      const equal = JSON.stringify(previousData) === JSON.stringify(currentData)
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

  return (
      servers.map((server) => (
      <tr ref={React.createRef()} key={category + server.planCode} className={`${changedRows.has(server.planCode) ? 'bg-yellow-200' : 'transition duration-1000 delay-150 even:bg-blue-300 odd:bg-blue-100'} font-mono`}>
        <td>{server.name}</td>
        <td>{server.cpu}</td>
        <td>{server.memory}</td>
        <td>{server.storage}</td>
        <td>{server.bandwidth}</td>
        <td>{server.price} {server.currencyCode}</td>
        <td className="flex flex-row justify-end">{server.status}<div className={statusColor.get(server.status)?.color + " pl-2"}><FontAwesomeIcon icon={statusColor.get(server.status)?.icon||faCircle} /></div></td>
        <td>{server.datacenters?.join(", ")||"-"}</td>
      </tr>
    ))
  )
}

function ServerCategories({ordered} : {ordered: {[key: string]: Server[]}}) {
  return (
    Object.entries(ordered).map(([category, servers]) => (
      <>
        <tr><td className="p-2 font-mono" colSpan={8}>{category||"Uncategorized"}</td></tr>
        <ServerLines category={category} servers={servers} />
      </>
    ))
  )
}

const ServersTable = ({data} : Props) => {
  if (!data || data.length === 0) {
    return
  }

  const serversByCategory = Object.groupBy(data, ( server: Server ) => server.category);

  const categoryOrder: {[key: string] :Server[]} = { "Kimsufi": [], "So you Start": [], "Rise": [], "": [], }
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
