import { useContext, useEffect, useRef, useState } from "react";
import GlobalContext from "../context/GlobalContext";

export const useWebSocket = (url) => {
  const [ws, setWs] = useState(null);
  const [reconnector, setReconnector] = useState(false);

  useEffect(() => {
    const socket = new WebSocket(url);
    setWs(socket);

    return () => {
      socket.close();
    };
  }, [url, reconnector]);

  const reconnect = () => {
    setReconnector((prev) => !prev);
  };

  return { ws, reconnect };
};
