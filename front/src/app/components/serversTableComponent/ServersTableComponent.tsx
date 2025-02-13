"use client";

import React from "react";

import ServersLineCategory from "../serversLineCategory/ServersLineCategory";

import { Server } from "../types";

interface Props {
  data: Server[];
}

// ServerLines displays the servers table lines

const ServersTable = ({ data }: Props) => {
  if (!data || data.length === 0) return <>Loading...</>;

  type Category = "Kimsufi" | "So you Start" | "Rise" | "uncategorized";
  type CategoryOrder = Record<Category, Server[]>;

  // Define the order of categories
  const categoryOrder: CategoryOrder = {
    Kimsufi: [],
    "So you Start": [],
    Rise: [],
    uncategorized: [],
  };

  // Group servers by category
  const serversByCategory = Object.groupBy(data, (server: Server) => server.category);

  // Merge the server and respect the categories order
  const ordered: CategoryOrder = { ...categoryOrder, ...serversByCategory };

  return (
    <table className="text-nowrap border-separate border-spacing-x-0 border-spacing-y-4">
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
        {Object.entries(ordered).map(([category, servers]) => (
          <ServersLineCategory key={category} category={category} servers={servers} />
        ))}
      </tbody>
    </table>
  );
};

export default ServersTable;
