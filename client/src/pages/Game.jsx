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
  const publisherSecret = "1234423hello";
  const user = "aidan";

  const x = 20;
  const y = 15;

  const bombs = 12;

  const [publicSecret, setPublicSecret] = useState("MY-PUBLIC-SECRET");

  const { ws: publisher, reconnect: reconnectPublisher } = useWebSocket(
    `${websocketServer}/publish?publisherSecret=${publisherSecret}&user=${user}&x=${x}&y=${y}&bombs=${bombs}&publicSecret=${publicSecret}`
  );

  const { ws: subscriber, reconnect: reconnectSubscriber } = useWebSocket(
    `${websocketServer}/subscribe?publicSecret=${publicSecret}`
  );

  const [board, setBoard] = useState([]);

  const [publisherIsOpen, setPublisherIsOpen] = useState(true);
  const [subscriberIsOpen, setSubscriberIsOpen] = useState(true);

  useEffect(() => {
    if (!publisher) return;

    publisher.onopen = () => {
      setPublisherIsOpen(true);
      console.log("WebSocket connected");
    };

    publisher.onmessage = (event) => {
      // Handle incoming messages here
    };

    publisher.onclose = () => {
      setPublisherIsOpen(false);
      console.log("publisher WebSocket disconnected");
      // Optional: Reconnect logic
    };

    publisher.onerror = (error) => {
      setPublisherIsOpen(false);
      console.error("WebSocket error:", error);
    };
  }, [publisher]);

  useEffect(() => {
    if (!subscriber) return;

    subscriber.onopen = () => {
      setSubscriberIsOpen(true);
      console.log("WebSocket connected");
    };

    subscriber.onmessage = (event) => {
      const board = JSON.parse(event.data);
      console.log("new board incoming");
      setBoard(board);
    };

    subscriber.onclose = () => {
      setSubscriberIsOpen(false);
      console.log("WebSocket disconnected");
      // Optional: Reconnect logic
    };

    subscriber.onerror = (error) => {
      setSubscriberIsOpen(false);
      console.error("WebSocket error:", error);
    };
  }, [subscriber]);

  const sendCoordinate = (x, y) => {
    // console.log(x, y);
    // send the coordinate
    const coordinates = {
      x: x,
      y: y,
    };
    console.log(coordinates);
    publisher.send(JSON.stringify(coordinates));
  };

  const handleContextMenu = (event) => {
    event.preventDefault();
  };

  return (
    <div onContextMenu={handleContextMenu}>
      {board.map((row, yCoord) => (
        <div key={yCoord}>
          {row.map((vox, xCoord) => (
            <Voxel
              key={xCoord}
              x={xCoord}
              y={yCoord}
              value={vox}
              sendCoordinate={sendCoordinate}
            />
          ))}
        </div>
      ))}
      game
      <button
        className={`btn btn-xs `}
        disabled={publisherIsOpen}
        onClick={() => {
          reconnectPublisher();
        }}
      >
        connect to publisher
      </button>
      <button
        className="btn btn-xs"
        disabled={subscriberIsOpen}
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
      <button
        className="btn btn-xs btn-primary"
        onClick={() => {
          sendCoordinate(0, 0);
        }}
      >
        send pub msg
      </button>
    </div>
  );
};

const Voxel = (props) => {
  const [flagged, setFlagged] = useState(false);

  const { x, y, value, sendCoordinate } = props;

  const rightClickAction = () => {
    setFlagged((prev) => !prev);
  };

  const leftClickAction = () => {
    if (!flagged) {
      sendCoordinate(x, y);
    }
  };

  return (
    <button
      onContextMenu={rightClickAction}
      className={`btn btn-xs btn-square ${
        value === "+" && !flagged && "text-transparent"
      } `}
      onClick={leftClickAction}
      disabled={value !== "+"}
    >
      {flagged ? "⚑" : <>{value === "+" ? "-" : value == "B" ? "💥" : value}</>}
    </button>
  );
};

export default Game;
