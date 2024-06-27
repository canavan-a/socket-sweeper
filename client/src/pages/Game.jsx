import { useContext, useEffect, useState } from "react";
import { json, useNavigate, useParams } from "react-router-dom";
import { useWebSocket } from "../hooks/WebsocketHooks";
import GlobalContext from "../context/GlobalContext";

const Game = () => {
  const userType = useParams();
  const navigate = useNavigate();

  const validUserTypes = ["sub", "pub"];

  const { websocketServer } = useContext(GlobalContext);

  //dynamically set this if we are making game or reconnecting to game
  const publisherSecret = "hello";
  const user = "aidan";

  const { ws: publisher, reconnect: reconnectPublisher } = useWebSocket(
    `${websocketServer}/publish?publisherSecret=${publisherSecret}&user=${user}`
  );

  const [board, setBoard] = useState([]);

  useEffect(() => {
    if (!validUserTypes.includes(userType.parameter)) {
      navigate("/");
    }
  }, []);

  useEffect(() => {
    if (!publisher) return;

    publisher.onopen = () => {
      console.log("WebSocket connected");
    };

    publisher.onmessage = (event) => {
      console.log("Message received:", event.data);
      // Handle incoming messages here
    };

    publisher.onclose = () => {
      console.log("WebSocket disconnected");
      // Optional: Reconnect logic
    };

    publisher.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }, [publisher]);

  const sendCoordinate = (x, y) => {
    console.log(x, y);
    // send the coordinate
    publisher.send();
  };

  return (
    <div>
      game
      <button className="btn btn-xs" onClick={reconnectPublisher}></button>
    </div>
  );
};

export default Game;
