import { useState } from "react";
import {
  Route,
  Routes,
  BrowserRouter as Router,
  Navigate,
} from "react-router-dom";
import Home from "./pages/Home";
import Game from "./pages/Game";
import { GlobalContextProvider } from "./context/GlobalContext";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <GlobalContextProvider>
        <Router>
          <Routes>
            <Route path="/*" element={<Home />} />
            <Route path="/game/:parameter" element={<Game />} />
          </Routes>
        </Router>
      </GlobalContextProvider>
    </div>
  );
}

export default App;
