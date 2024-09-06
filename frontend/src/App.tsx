import { BrowserRouter, Route, Routes } from "react-router-dom";
import Nav from "./components/Nav";
import Search from "./views/Search";

function App() {
  return (
    <>
      <BrowserRouter>
        <Nav loggedIn={false} />
        <Routes>
          <Route
            path="/"
            element={<Search />}
          />
          <Route
            path="/weather"
            element={<div />}
          />
          <Route
            path="/register"
            element={<div />}
          />
          <Route
            path="/login"
            element={<div />}
          />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
