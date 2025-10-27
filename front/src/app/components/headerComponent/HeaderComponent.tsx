import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTelegram } from "@fortawesome/free-brands-svg-icons";
import Status from "../statusComponent/StatusComponent";
import { StatusProps } from "../statusComponent/type";

export default function HeaderComponent({
  error,
  data,
  connectionRestored,
  setConnectionRestored,
  lastMessage,
}: StatusProps) {
  return (
    <div className="flex justify-evenly w-full p-2">
      <div className="flex flex-col justify-center">
        <div>Receive notifications on Telegram</div>
        <a className="flex text-blue-400 space-x-1" href="https://t.me/KimsufiNotifierBot">
          <div>
            <FontAwesomeIcon icon={faTelegram} />
          </div>
          <div>t.me/KimsufiNotifierBot</div>
        </a>
      </div>
      <h1 className="text-center text-xl font-bold">OVH Eco server availability</h1>
      <Status
        error={error}
        data={data}
        connectionRestored={connectionRestored}
        setConnectionRestored={setConnectionRestored}
        lastMessage={lastMessage}
      />
    </div>
  );
}
