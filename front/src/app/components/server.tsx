import { ReactNode } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle, faTimesCircle } from '@fortawesome/free-regular-svg-icons';

interface Props {
    data: Object
}

function ServerLines({ servers }) {
  let statusColor = new Map<string, Object>([
      ["available", {"color":"text-lime-600", "icon":faCheckCircle}],
      ["unavailable", {"color":"text-rose-600", "icon":faTimesCircle}],
  ]);

  return (
    servers.map((server) => (
      <tr key={server.planCode} className="font-mono even:bg-blue-300 odd:bg-blue-100">
        <td>{server.planCode}</td>
        <td>{server.category||"uncategorized"}</td>
        <td>{server.name}</td>
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

  const serversbycategory = Object.groupBy(data, ({ category }) => category);

  let categoryorder = { "Kimsufi": [], "So you Start": [], "Rise": [], "": [], }
  const ordered = Object.assign(categoryorder, serversbycategory);

  const tableBody = Object.entries(ordered).map(([category, servers]) => (
    <>
      <ServerLines key={category + " servers"} servers={servers} />
      <tr key={category + " separator0"}><td className="p-2" colSpan={6}></td></tr>
      <tr key={category + " separator1"}><td className="p-2" colSpan={6}></td></tr>
    </>
  ));

  return (
      <table className="text-nowrap">
        <thead>
          <tr>
            <th className="p-4">Plan Code</th>
            <th className="p-4">Category</th>
            <th className="p-4">Name</th>
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
