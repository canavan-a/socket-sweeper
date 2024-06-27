import { useNavigate } from "react-router-dom";

const Home = () => {
  const navigate = useNavigate();
  const goToGame = () => {
    navigate("/game/pub");
  };
  return (
    <div>
      home
      <button className="btn btn-xs" onClick={goToGame}>
        navigate to game
      </button>
    </div>
  );
};

export default Home;
