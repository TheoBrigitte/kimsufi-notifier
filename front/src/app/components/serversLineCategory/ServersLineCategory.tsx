import React from "react";

import { ServerLineData } from "../serverLineData/ServerLineData";

import { Server } from "../types";

// ServerCategories displays the servers table lines with a category line before
export default function ServersLineCategory({
  category,
  servers,
}: {
  category: string;
  servers: Server[];
}) {
  return (
    <>
      <tr>
        <td className="p-4 font-mono font-semibold border-b-2 border-gray-300" colSpan={8}>
          {category}
        </td>
      </tr>
      {servers.map((server) => (
        <ServerLineData key={category + server.planCode} server={server} />
      ))}
    </>
  );
}
