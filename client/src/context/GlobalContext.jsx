import { createContext, useState, useEffect } from "react";

const GlobalContext = createContext();

export const GlobalContextProvider = ({ children }) => {
  const envUrl = import.meta.env.VITE_APP_SERVER_URI;
  const [server, setServer] = useState(
    envUrl.includes("localhost") ? envUrl : ""
  );

  const envWsUrl = import.meta.env.VITE_APP_WS_ENDPOINT;

  const [websocketServer, setWebsocketServer] = useState(
    envWsUrl.includes("localhost")
      ? envWsUrl
      : `wss://${window.location.hostname}`
  );

  const contextValue = {
    server,
    websocketServer,
  };

  return (
    <GlobalContext.Provider value={contextValue}>
      {children}
    </GlobalContext.Provider>
  );
};

export default GlobalContext;
