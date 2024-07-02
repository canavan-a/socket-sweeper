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
  const publisherSecret = "hellosasasass";
  const user = "aidan";

  const x = 50;
  const y = 75;

  const bombs = 20;

  const [publicSecret, setPublicSecret] = useState("MY-PUBLIC-SECRET");

  const { ws: publisher, reconnect: reconnectPublisher } = useWebSocket(
    `${websocketServer}/publish?publisherSecret=${publisherSecret}&user=${user}&x=${x}&y=${y}&bombs=${bombs}&publisherSecret=${publicSecret}`
  );

  const { ws: subscriber, reconnect: reconnectSubscriber } = useWebSocket(
    `${websocketServer}/suscribe?publicSecret=${publicSecret}`
  );

  const [board, setBoard] = useState([]);

  // useEffect(() => {
  //   if (!validUserTypes.includes(userType.parameter)) {
  //     navigate("/");
  //   }
  // }, []);

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
      console.log("publisher WebSocket disconnected");
      // Optional: Reconnect logic
    };

    publisher.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }, [publisher]);

  useEffect(() => {
    if (!subscriber) return;

    subscriber.onopen = () => {
      console.log("WebSocket connected");
    };

    subscriber.onmessage = (event) => {
      console.log(event.data);
      console.log("Message received:", event.data);
      // Handle incoming messages here
    };

    subscriber.onclose = () => {
      console.log("WebSocket disconnected");
      // Optional: Reconnect logic
    };

    subscriber.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }, [subscriber]);

  const sendCoordinate = (x, y) => {
    console.log(x, y);
    // send the coordinate
    publisher.send();
  };

  return (
    <div>
      game
      <button
        className="btn btn-xs"
        onClick={() => {
          reconnectPublisher();
        }}
      >
        connect to publisher
      </button>
      <button
        className="btn btn-xs"
        onClick={() => {
          reconnectSubscriber();
        }}
      >
        connect to subscriber
      </button>
      <input
        className="input input-xs input-bordered"
        value={publicSecret}
        onChange={(e) => {
          setPublicSecret(e.target.value);
        }}
      ></input>
    </div>
  );
};

export default Game;
