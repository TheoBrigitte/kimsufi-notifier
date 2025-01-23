import { ServerLineData } from "../serverLineData/ServerLineData";

import { Server } from "../types";

interface Props {
  category: string;
  servers: Server[];
}

// ServerCategories displays the servers table lines with a category line before
export default function ServersLineCategory({ category, servers }: Props) {
  return (
    <>
      <tr>
        <td className="p-4 font-mono font-semibold border-b-2 border-gray-300" colSpan={8}>
          {category}
        </td>
      </tr>
      {servers.some((server) => server) ? (
        servers.map((server) => <ServerLineData key={category + server.planCode} server={server} />)
      ) : (
        <tr>
          <td colSpan={8} className="text-center text-gray-600">
            No servers yet
          </td>
        </tr>
      )}
    </>
  );
}
