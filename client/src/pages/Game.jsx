import { useContext, useEffect, useState } from "react";
import { json, useNavigate, useParams } from "react-router-dom";
import { useWebSocket } from "../hooks/WebsocketHooks";
import GlobalContext from "../context/GlobalContext";

const Game = () => {
  const userType = useParams();
  const navigate = useNavigate();

  const validUserTypes = ["sub", "pub"];

  const { websocketServer } = useContext(GlobalContext);

  const { ws, reconnect } = useWebSocket(`${websocketServer}/`);

  const [board, setBoard] = useState([]);

  useEffect(() => {
    if (!validUserTypes.includes(userType.parameter)) {
      navigate("/");
    }
  }, []);

  useEffect(() => {
    if (!ws) return;

    ws.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      console.log("Message received:", event.data);
      // Handle incoming messages here
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
      // Optional: Reconnect logic
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }, [ws]);

  const sendCoordinate = (x, y) => {
    console.log(x, y);
    // send the coordinate
    ws.send();
  };

  return (
    <div>
      game
      <button className="btn btn-xs" onClick={reconnect}></button>
    </div>
  );
};

export default Game;
