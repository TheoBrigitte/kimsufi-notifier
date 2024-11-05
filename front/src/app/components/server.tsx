import { ReactNode } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle, faTimesCircle } from '@fortawesome/free-regular-svg-icons';

interface Props {
    servers: Array
}

const ServerLine = ({servers} : Props) => {
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
  );
};

export default ServerLine;
