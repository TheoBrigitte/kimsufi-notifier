import ServersLineCategory from "../serversLineCategory/ServersLineCategory";

import { Server } from "../types";

interface Props {
  data: Server[];
}

const columnsHead = [
  "Name",
  "CPU",
  "RAM",
  "Storage",
  "Bandwidth",
  "Price",
  "Status",
  "Datacenters",
];

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

  //checking uncategorized servers
  const checkedUncategorized = data.map((server) =>
    !server.category ? { ...server, category: "uncategorized" } : server
  );

  // Group servers by category
  const serversByCategory = Object.groupBy(
    checkedUncategorized,
    (server: Server) => server.category
  );

  // Merge the server and respect the categories order
  const ordered: CategoryOrder = { ...categoryOrder, ...serversByCategory };

  return (
    <table className="text-nowrap border-separate border-spacing-x-0">
      <thead>
        <tr>
          {columnsHead.map((columnHead) => (
            <th key={columnHead} className="p-2">
              {columnHead}
            </th>
          ))}
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
