import { createContext, useState } from "react";

const GlobalContext = createContext();

export const GlobalContextProvider = ({ children }) => {
  const [server, setServer] = useState(
    `${import.meta.env.VITE_APP_SERVER_URI}`
  );

  const [websocketServer, setWebsocketServer] = useState(
    `${import.meta.env.VITE_APP_WS_ENDPOINT}`
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
